#!/usr/bin/env bash
#
# benchmark.sh — Compare synche vs rsync vs scp for two scenarios:
#   Case 1: Upload 100MB of data that doesn't exist on the server
#   Case 2: Upload 100MB of data where 50MB already exists on the server
#
# scp is included in remote mode to provide a baseline for raw SSH transfer
# speed (no checksumming or delta logic). Use --no-scp to disable.
#
# This script is fully repeatable: deterministic data, clean state each run.
#
# Usage:
#   # Local mode (starts servers automatically):
#   ./benchmark.sh
#
#   # Remote mode with rsync daemon (no encryption, fair protocol comparison):
#   ./benchmark.sh --server YOUR_SERVER_IP --rsync-mode daemon
#
#   # Remote mode with rsync over SSH (real-world comparison):
#   ./benchmark.sh --server YOUR_SERVER_IP --rsync-mode ssh --ssh-user root
#
#   # Both rsync modes:
#   ./benchmark.sh --server YOUR_SERVER_IP --rsync-mode both --ssh-user root
#
# See --setup-remote for server setup instructions.
#
set -euo pipefail

# ---------------------------------------------------------------------------
# Configuration (defaults)
# ---------------------------------------------------------------------------
ROUNDS=3
SEED=42
SECOND_SEED=99
SIZE_MB=100
HALF_MB=50
SERVER_HOST=""          # empty = local mode
SYNCHE_PORT=8420
RSYNC_PORT=8730
RSYNC_MODE="daemon"    # daemon, ssh, or both
CONCURRENCY=8
BENCH_DIR="/tmp/synche-bench-$$"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
SYNCHE_CLIENT="${SCRIPT_DIR}/bin/synche-client"
SYNCHE_SERVER="${SCRIPT_DIR}/bin/synche-server"
BENCHDATA="${SCRIPT_DIR}/bin/benchdata"
SETUP_REMOTE=false
SSH_USER=""
RSYNC_SSH_PATH="/tmp/rsync-bench"  # remote path for rsync over SSH
SCP_PATH="/tmp/scp-bench"          # remote path for scp uploads
SCP_ENABLED=true                   # scp benchmark (remote mode only)
SYNCHE_API_KEY=""                   # API key for synche server (optional)
SYNCHE_STORE="/var/lib/synche-store" # remote synche store directory

# ---------------------------------------------------------------------------
# Parse args
# ---------------------------------------------------------------------------
while [[ $# -gt 0 ]]; do
    case "$1" in
        --server)        SERVER_HOST="$2"; shift 2;;
        --rounds)        ROUNDS="$2"; shift 2;;
        --synche-port)   SYNCHE_PORT="$2"; shift 2;;
        --rsync-port)    RSYNC_PORT="$2"; shift 2;;
        --rsync-mode)    RSYNC_MODE="$2"; shift 2;;
        --rsync-path)    RSYNC_SSH_PATH="$2"; shift 2;;
        --concurrency)   CONCURRENCY="$2"; shift 2;;
        --size)          SIZE_MB="$2"; HALF_MB=$(($2 / 2)); shift 2;;
        --setup-remote)  SETUP_REMOTE=true; shift;;
        --ssh-user)      SSH_USER="$2"; shift 2;;
        --scp-path)      SCP_PATH="$2"; shift 2;;
        --no-scp)        SCP_ENABLED=false; shift;;
        --api-key)       SYNCHE_API_KEY="$2"; shift 2;;
        --synche-store)  SYNCHE_STORE="$2"; shift 2;;
        --help|-h)
            cat <<EOF
Usage: $0 [OPTIONS]

Options:
  --server HOST      Remote server IP/hostname (omit for local mode)
  --rounds N         Number of rounds per case (default: 3)
  --synche-port N    Synche server port (default: 8420)
  --rsync-port N     Rsync daemon port (default: 8730)
  --rsync-mode MODE  "daemon" (no encryption), "ssh", or "both" (default: daemon)
  --rsync-path PATH  Remote directory for rsync over SSH (default: /tmp/rsync-bench)
  --concurrency N    Synche upload concurrency (default: 8)
  --size N           Total data size in MiB (default: 100)
  --ssh-user USER    SSH user for remote rsync over SSH (default: current user)
  --scp-path PATH    Remote directory for scp uploads (default: /tmp/scp-bench)
  --no-scp           Disable scp benchmark
  --api-key KEY      API key for synche server authentication
  --synche-store DIR Remote synche store directory (default: /var/lib/synche-store)
  --setup-remote     Print remote server setup instructions and exit
  --help             Show this help
EOF
            exit 0
            ;;
        *) echo "Unknown arg: $1 (use --help)"; exit 1;;
    esac
done

# Validate rsync mode.
case "$RSYNC_MODE" in
    daemon|ssh|both) ;;
    *) echo "Invalid --rsync-mode: $RSYNC_MODE (must be daemon, ssh, or both)"; exit 1;;
esac

# ---------------------------------------------------------------------------
# Print remote setup instructions
# ---------------------------------------------------------------------------
if $SETUP_REMOTE; then
    cat <<'SETUP'
=== Remote Server Setup ===

1. Copy the synche-server binary to your server:

   scp bin/synche-server user@YOUR_SERVER:/usr/local/bin/

2. Start the synche server:

   synche-server --addr :8420 --store /tmp/synche-store


--- For --rsync-mode daemon (recommended for fair comparison) ---

3. Create /tmp/rsyncd.conf:

   pid file = /tmp/rsyncd.pid
   port = 8730
   use chroot = no
   read only = no
   uid = nobody
   gid = nogroup

   [bench]
       path = /tmp/rsync-bench
       read only = no

4. Create the directory and start the daemon:

   mkdir -p -m 777 /tmp/rsync-bench
   rsync --daemon --config=/tmp/rsyncd.conf --no-detach


--- For --rsync-mode ssh ---

3. Ensure SSH access to the server and create the directory:

   ssh user@YOUR_SERVER 'mkdir -p /tmp/rsync-bench'


--- Run the benchmark ---

   # Daemon mode (no encryption — fair protocol comparison):
   ./benchmark.sh --server YOUR_SERVER_IP --rsync-mode daemon

   # SSH mode (real-world comparison, includes encryption overhead):
   ./benchmark.sh --server YOUR_SERVER_IP --rsync-mode ssh --ssh-user root

   # Both modes:
   ./benchmark.sh --server YOUR_SERVER_IP --rsync-mode both --ssh-user root

SETUP
    exit 0
fi

# ---------------------------------------------------------------------------
# Determine mode
# ---------------------------------------------------------------------------
if [[ -z "$SERVER_HOST" ]]; then
    MODE="local"
    SYNCHE_URL="http://localhost:${SYNCHE_PORT}"
    RSYNC_HOST="localhost"
    # Local mode always uses daemon (no SSH to localhost).
    if [[ "$RSYNC_MODE" == "ssh" ]]; then
        echo "[warning] --rsync-mode ssh not supported in local mode, using daemon"
        RSYNC_MODE="daemon"
    fi
    if [[ "$RSYNC_MODE" == "both" ]]; then
        echo "[warning] --rsync-mode both not supported in local mode, using daemon"
        RSYNC_MODE="daemon"
    fi
    if $SCP_ENABLED; then
        echo "[warning] scp benchmark not supported in local mode, disabling"
        SCP_ENABLED=false
    fi
else
    MODE="remote"
    SYNCHE_URL="http://${SERVER_HOST}:${SYNCHE_PORT}"
    RSYNC_HOST="${SERVER_HOST}"
fi

# Build SSH target.
if [[ -n "$SSH_USER" ]]; then
    RSYNC_SSH_TARGET="${SSH_USER}@${RSYNC_HOST}"
    SCP_SSH_TARGET="${SSH_USER}@${RSYNC_HOST}"
    SYNCHE_SSH_TARGET="${SSH_USER}@${RSYNC_HOST}"
else
    RSYNC_SSH_TARGET="${RSYNC_HOST}"
    SCP_SSH_TARGET="${RSYNC_HOST}"
    SYNCHE_SSH_TARGET="${RSYNC_HOST}"
fi

# Build synche client API key flag.
SYNCHE_API_KEY_FLAG=""
if [[ -n "$SYNCHE_API_KEY" ]]; then
    SYNCHE_API_KEY_FLAG="--api-key ${SYNCHE_API_KEY}"
fi

# Determine which rsync modes to run.
declare -a RSYNC_MODES
case "$RSYNC_MODE" in
    daemon) RSYNC_MODES=(daemon);;
    ssh)    RSYNC_MODES=(ssh);;
    both)   RSYNC_MODES=(daemon ssh);;
esac

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
SYNCHE_PID=0
RSYNC_PID=0

cleanup() {
    echo ""
    echo "[cleanup] stopping background processes..."
    if [[ "$MODE" == "local" ]]; then
        kill "$SYNCHE_PID" 2>/dev/null || true
        kill "$RSYNC_PID" 2>/dev/null || true
        wait "$SYNCHE_PID" 2>/dev/null || true
        wait "$RSYNC_PID" 2>/dev/null || true
    fi
    rm -rf "$BENCH_DIR"
    echo "[cleanup] done"
}
trap cleanup EXIT

start_synche_server() {
    if [[ "$MODE" == "local" ]]; then
        rm -rf "${BENCH_DIR}/synche-store"
        local api_key_flag=""
        if [[ -n "$SYNCHE_API_KEY" ]]; then
            api_key_flag="--api-key ${SYNCHE_API_KEY}"
        fi
        "$SYNCHE_SERVER" --addr ":${SYNCHE_PORT}" --store "${BENCH_DIR}/synche-store" \
            $api_key_flag >/dev/null 2>&1 &
        SYNCHE_PID=$!
        sleep 0.5
    fi
}

reset_synche_server() {
    if [[ "$MODE" == "local" ]]; then
        # Local: stop, wipe store, restart.
        kill "$SYNCHE_PID" 2>/dev/null || true
        wait "$SYNCHE_PID" 2>/dev/null || true
        rm -rf "${BENCH_DIR}/synche-store"
        local api_key_flag=""
        if [[ -n "$SYNCHE_API_KEY" ]]; then
            api_key_flag="--api-key ${SYNCHE_API_KEY}"
        fi
        "$SYNCHE_SERVER" --addr ":${SYNCHE_PORT}" --store "${BENCH_DIR}/synche-store" \
            $api_key_flag >/dev/null 2>&1 &
        SYNCHE_PID=$!
        sleep 0.5
    fi
}

# delete_synche_manifest <manifest_id>
# Deletes a manifest and its orphaned chunks via the API.
delete_synche_manifest() {
    local id="$1"
    if [[ -z "$id" ]]; then
        return
    fi
    local auth_header=""
    if [[ -n "$SYNCHE_API_KEY" ]]; then
        auth_header="Authorization: Bearer ${SYNCHE_API_KEY}"
    fi
    if ! curl -sf --max-time 5 -X DELETE "${SYNCHE_URL}/api/manifest/${id}" \
        ${auth_header:+-H "$auth_header"} -o /dev/null; then
        : # silently ignore (may be duplicate or already deleted)
    fi
}

# run_synche <source_file> [extra_args...]
# Runs the synche client and prints the manifest ID to stdout.
run_synche() {
    local src="$1"
    shift
    local output
    output=$("$SYNCHE_CLIENT" \
        --source "$src" \
        --server "$SYNCHE_URL" \
        --cache-dir "${BENCH_DIR}/synche-cache" \
        --concurrency "$CONCURRENCY" \
        $SYNCHE_API_KEY_FLAG "$@" 2>&1)
    local mid
    mid=$(echo "$output" | grep -oP 'manifest ID:\s+\K\S+' || true)
    echo "$mid"
}

# run_synche_dir <source_dir> [extra_args...]
# Runs the synche client on a directory and prints all manifest IDs (one per line).
run_synche_dir() {
    local src="$1"
    shift
    local output
    output=$("$SYNCHE_CLIENT" \
        --source "$src" \
        --server "$SYNCHE_URL" \
        --cache-dir "${BENCH_DIR}/synche-cache" \
        --concurrency "$CONCURRENCY" \
        $SYNCHE_API_KEY_FLAG "$@" 2>&1)
    echo "$output" | grep -oP 'manifest ID:\s+\K\S+' | sort -u
}

# delete_synche_manifests <manifest_ids...>
# Deletes multiple manifests (one ID per argument or newline-separated).
delete_synche_manifests() {
    for id in "$@"; do
        # Handle newline-separated IDs from run_synche_dir.
        while IFS= read -r mid; do
            delete_synche_manifest "$mid"
        done <<< "$id"
    done
}

# rsync_send_dir <mode> <src_dir> [extra flags...]
# Sends a directory via rsync using the specified mode.
rsync_send_dir() {
    local rmode="$1" src="$2"
    shift 2
    local extra_flags=("$@")

    case "$rmode" in
        daemon)
            rsync "${extra_flags[@]}" -r "$src/" \
                "rsync://${RSYNC_HOST}:${RSYNC_PORT}/bench/"
            ;;
        ssh)
            rsync "${extra_flags[@]}" -r "$src/" \
                "${RSYNC_SSH_TARGET}:${RSYNC_SSH_PATH}/"
            ;;
    esac
}

# scp_send_dir <src_dir>
# Sends a directory via scp.
scp_send_dir() {
    local src="$1"
    scp -q -r "$src/" "${SCP_SSH_TARGET}:${SCP_PATH}/"
}

start_rsync_daemon() {
    if [[ "$MODE" == "local" ]]; then
        local module_path="$1"
        mkdir -p "$module_path"
        chmod 777 "$module_path"
        cat > "${BENCH_DIR}/rsyncd.conf" <<CONF
pid file = ${BENCH_DIR}/rsyncd.pid
port = ${RSYNC_PORT}
use chroot = no
read only = no

[bench]
    path = ${module_path}
    read only = no
CONF
        rsync --daemon --config="${BENCH_DIR}/rsyncd.conf" --no-detach \
            >/dev/null 2>&1 &
        RSYNC_PID=$!
        sleep 0.5
    fi
}

stop_rsync_daemon() {
    if [[ "$MODE" == "local" ]]; then
        kill "$RSYNC_PID" 2>/dev/null || true
        wait "$RSYNC_PID" 2>/dev/null || true
        RSYNC_PID=0
    fi
}

# rsync_send <mode> <src_file> <dest_name> [extra flags...]
# Sends a file via rsync using the specified mode.
# Pass --quiet as an extra flag to suppress progress (used for pre-seeding).
rsync_send() {
    local rmode="$1" src="$2" dest_name="$3"
    shift 3
    local extra_flags=("$@")

    # Add --progress unless --quiet was passed.
    local show_progress=true
    for f in "${extra_flags[@]}"; do
        if [[ "$f" == "--quiet" ]]; then
            show_progress=false
            break
        fi
    done
    if $show_progress; then
        extra_flags+=(--progress)
    fi

    case "$rmode" in
        daemon)
            rsync "${extra_flags[@]}" "$src" \
                "rsync://${RSYNC_HOST}:${RSYNC_PORT}/bench/${dest_name}"
            ;;
        ssh)
            rsync "${extra_flags[@]}" "$src" \
                "${RSYNC_SSH_TARGET}:${RSYNC_SSH_PATH}/${dest_name}"
            ;;
    esac
}

# rsync_reset <mode>
# Clears all files in the rsync destination.
rsync_reset() {
    local rmode="$1"
    local empty_dir="${BENCH_DIR}/empty"
    mkdir -p "$empty_dir"

    case "$rmode" in
        daemon)
            rsync --delete -r "${empty_dir}/" \
                "rsync://${RSYNC_HOST}:${RSYNC_PORT}/bench/" 2>/dev/null || true
            ;;
        ssh)
            ssh "${RSYNC_SSH_TARGET}" "rm -rf ${RSYNC_SSH_PATH}/* ${RSYNC_SSH_PATH}/.[!.]*" \
                2>/dev/null || true
            ;;
    esac
}

rsync_mode_label() {
    case "$1" in
        daemon) echo "rsync(daemon)";;
        ssh)    echo "rsync(ssh)";;
    esac
}

# scp_reset — remove all files in the scp destination directory.
scp_reset() {
    ssh "${SCP_SSH_TARGET}" "rm -rf ${SCP_PATH}/* ${SCP_PATH}/.[!.]*" 2>/dev/null || true
}

# scp_send <src_file> <dest_name>
scp_send() {
    local src="$1" dest_name="$2"
    scp -q "$src" "${SCP_SSH_TARGET}:${SCP_PATH}/${dest_name}"
}

check_server() {
    local label="$1"
    local url="$2"
    if ! curl -sf "$url" >/dev/null 2>&1; then
        echo "[error] ${label} is not reachable at ${url}"
        echo "        Make sure the server is running. Use --setup-remote for instructions."
        exit 1
    fi
}

# ---------------------------------------------------------------------------
# Setup
# ---------------------------------------------------------------------------
echo "================================================================="
echo "  synche vs rsync vs scp benchmark"
echo "================================================================="
echo "  mode:         ${MODE}"
if [[ "$MODE" == "remote" ]]; then
echo "  server:       ${SERVER_HOST}"
fi
echo "  synche url:   ${SYNCHE_URL}"
echo "  rsync mode:   ${RSYNC_MODE}"
if [[ " ${RSYNC_MODES[*]} " =~ " daemon " ]]; then
echo "  rsync daemon: ${RSYNC_HOST}:${RSYNC_PORT}"
fi
if [[ " ${RSYNC_MODES[*]} " =~ " ssh " ]]; then
echo "  rsync ssh:    ${RSYNC_SSH_TARGET}:${RSYNC_SSH_PATH}"
fi
if $SCP_ENABLED; then
echo "  scp:          ${SCP_SSH_TARGET}:${SCP_PATH}"
fi
if [[ -n "$SYNCHE_API_KEY" ]]; then
echo "  api key:      enabled"
fi
echo "  rounds:       ${ROUNDS}"
echo "  data size:    ${SIZE_MB} MiB"
echo "  concurrency:  ${CONCURRENCY}"
echo "  seed:         ${SEED}"
echo "  temp dir:     ${BENCH_DIR}"
echo "================================================================="
echo ""

mkdir -p "$BENCH_DIR"

# Start / verify servers.
if [[ "$MODE" == "local" ]]; then
    mkdir -p "${BENCH_DIR}/rsync-dest"
    start_synche_server
    start_rsync_daemon "${BENCH_DIR}/rsync-dest"

    sleep 0.3
    check_server "synche (local)" "${SYNCHE_URL}/api/health"
    echo "[preflight] local synche server: OK"
    if ! rsync "rsync://localhost:${RSYNC_PORT}/" >/dev/null 2>&1; then
        echo "[error] local rsync daemon failed to start"
        exit 1
    fi
    echo "[preflight] local rsync daemon: OK"
    echo ""
else
    echo "[preflight] checking remote servers..."
    check_server "synche" "${SYNCHE_URL}/api/health"
    echo "  synche server: OK"

    if [[ " ${RSYNC_MODES[*]} " =~ " daemon " ]]; then
        if ! rsync "rsync://${RSYNC_HOST}:${RSYNC_PORT}/" >/dev/null 2>&1; then
            echo "  [error] rsync daemon not reachable at ${RSYNC_HOST}:${RSYNC_PORT}"
            echo "          Use --setup-remote for instructions."
            exit 1
        fi
        echo "  rsync daemon:  OK"
    fi

    if [[ " ${RSYNC_MODES[*]} " =~ " ssh " ]]; then
        if ! ssh -o ConnectTimeout=5 "${RSYNC_SSH_TARGET}" "echo ok" >/dev/null 2>&1; then
            echo "  [error] SSH not reachable at ${RSYNC_SSH_TARGET}"
            echo "          Check SSH access and use --ssh-user if needed."
            exit 1
        fi
        # Ensure remote directory exists.
        ssh "${RSYNC_SSH_TARGET}" "mkdir -p ${RSYNC_SSH_PATH}"
        echo "  rsync ssh:     OK"
    fi

    if $SCP_ENABLED; then
        # SSH was already validated above if rsync ssh mode is active,
        # but check again in case only scp is used.
        if ! ssh -o ConnectTimeout=5 "${SCP_SSH_TARGET}" "echo ok" >/dev/null 2>&1; then
            echo "  [error] SSH not reachable at ${SCP_SSH_TARGET}"
            echo "          Check SSH access and use --ssh-user if needed."
            exit 1
        fi
        ssh "${SCP_SSH_TARGET}" "mkdir -p ${SCP_PATH}"
        echo "  scp:           OK"
    fi
    echo ""
fi

# Generate deterministic test data.
echo "[setup] generating ${SIZE_MB} MiB test file (seed=${SEED})..."
"$BENCHDATA" --seed "$SEED" --size "$SIZE_MB" --out "${BENCH_DIR}/source-full"

echo "[setup] generating ${HALF_MB} MiB pre-existing data (seed=${SECOND_SEED})..."
"$BENCHDATA" --seed "$SECOND_SEED" --size "$HALF_MB" --out "${BENCH_DIR}/preexist-half"

# Case-2 source: first half = preexist data, second half = new data.
echo "[setup] building case-2 source file (${HALF_MB}MiB existing + ${HALF_MB}MiB new)..."
cp "${BENCH_DIR}/preexist-half" "${BENCH_DIR}/source-case2"
dd if="${BENCH_DIR}/source-full" bs=1M count="$HALF_MB" >> "${BENCH_DIR}/source-case2" 2>/dev/null

# Case-2 "previous version" for rsync: first half matches, second half is different.
echo "[setup] building case-2 rsync previous-version file..."
cp "${BENCH_DIR}/preexist-half" "${BENCH_DIR}/rsync-prev-case2"
dd if=/dev/zero bs=1M count="$HALF_MB" >> "${BENCH_DIR}/rsync-prev-case2" 2>/dev/null

# Case 3/4: simulated disk image — 100 files of 0.9 MiB each (~90 MiB total).
DISK_FILES=100
DISK_FILE_SIZE_KB=921  # 0.9 MiB ~ 921 KB
DISK_CHANGED=10        # number of files that change between versions
DISK_TOTAL_KB=$((DISK_FILES * DISK_FILE_SIZE_KB))
DISK_TOTAL_MB=$(( (DISK_TOTAL_KB + 1023) / 1024 ))  # ceiling
DISK_CHANGED_KB=$((DISK_CHANGED * DISK_FILE_SIZE_KB))

echo "[setup] building ${DISK_TOTAL_MB} MiB disk image (${DISK_FILES} x ${DISK_FILE_SIZE_KB}KB files)..."
"$BENCHDATA" --seed 7777 --size "$DISK_TOTAL_MB" --out "${BENCH_DIR}/disk-original"
truncate -s "${DISK_TOTAL_KB}K" "${BENCH_DIR}/disk-original"

echo "[setup] building modified disk image (${DISK_CHANGED}/${DISK_FILES} files changed, ~${DISK_CHANGED_KB}KB delta)..."
cp "${BENCH_DIR}/disk-original" "${BENCH_DIR}/disk-modified"
# Overwrite 10 file-sized regions at evenly spaced positions with different data.
for ((i=0; i<DISK_CHANGED; i++)); do
    file_idx=$(( i * (DISK_FILES / DISK_CHANGED) ))  # files 0, 10, 20, ...
    offset_kb=$(( file_idx * DISK_FILE_SIZE_KB ))
    "$BENCHDATA" --seed $((9000 + i)) --size 1 --out "${BENCH_DIR}/disk-patch"
    truncate -s "${DISK_FILE_SIZE_KB}K" "${BENCH_DIR}/disk-patch"
    dd if="${BENCH_DIR}/disk-patch" of="${BENCH_DIR}/disk-modified" \
        bs=1024 seek="$offset_kb" conv=notrunc 2>/dev/null
done
rm -f "${BENCH_DIR}/disk-patch"

# Case-4 rsync "previous version" = disk-original (rsync will delta-sync to disk-modified).
cp "${BENCH_DIR}/disk-original" "${BENCH_DIR}/rsync-prev-case4"

# Case 5/6: directory of 100 individual files (0.9 MiB each).
echo "[setup] building directory with ${DISK_FILES} x ${DISK_FILE_SIZE_KB}KB files..."
mkdir -p "${BENCH_DIR}/dir-original"
for ((i=0; i<DISK_FILES; i++)); do
    fname=$(printf "file_%03d.bin" "$i")
    "$BENCHDATA" --seed $((2000 + i)) --size 1 --out "${BENCH_DIR}/dir-original/${fname}"
    truncate -s "${DISK_FILE_SIZE_KB}K" "${BENCH_DIR}/dir-original/${fname}"
done

echo "[setup] building modified directory (${DISK_CHANGED}/${DISK_FILES} files replaced)..."
cp -r "${BENCH_DIR}/dir-original" "${BENCH_DIR}/dir-modified"
for ((i=0; i<DISK_CHANGED; i++)); do
    file_idx=$(( i * (DISK_FILES / DISK_CHANGED) ))
    fname=$(printf "file_%03d.bin" "$file_idx")
    "$BENCHDATA" --seed $((5000 + i)) --size 1 --out "${BENCH_DIR}/dir-modified/${fname}"
    truncate -s "${DISK_FILE_SIZE_KB}K" "${BENCH_DIR}/dir-modified/${fname}"
done

echo "[setup] checksums:"
md5sum "${BENCH_DIR}/source-full" "${BENCH_DIR}/source-case2" "${BENCH_DIR}/preexist-half" \
    "${BENCH_DIR}/disk-original" "${BENCH_DIR}/disk-modified"
echo ""

# ---------------------------------------------------------------------------
# Results storage — one array per (case, tool+mode)
# ---------------------------------------------------------------------------
declare -a CASE1_SYNCHE CASE2_SYNCHE CASE3_SYNCHE CASE4_SYNCHE CASE5_SYNCHE CASE6_SYNCHE
declare -a CASE1_RSYNC_DAEMON CASE1_RSYNC_SSH
declare -a CASE2_RSYNC_DAEMON CASE2_RSYNC_SSH
declare -a CASE3_RSYNC_DAEMON CASE3_RSYNC_SSH
declare -a CASE4_RSYNC_DAEMON CASE4_RSYNC_SSH
declare -a CASE5_RSYNC_DAEMON CASE5_RSYNC_SSH
declare -a CASE6_RSYNC_DAEMON CASE6_RSYNC_SSH
declare -a CASE1_SCP CASE2_SCP CASE3_SCP CASE4_SCP CASE5_SCP CASE6_SCP

# ---------------------------------------------------------------------------
# Case 1: Fresh upload (nothing on server)
# ---------------------------------------------------------------------------
echo "================================================================="
echo "  CASE 1: Fresh ${SIZE_MB} MiB upload (no data on server)"
echo "================================================================="
echo ""

# In local mode, reset once for a clean slate.
if [[ "$MODE" == "local" ]]; then
    reset_synche_server
fi

for ((r=1; r<=ROUNDS; r++)); do
    echo "--- Round $r/$ROUNDS ---"

    # -- synche --
    rm -rf "${BENCH_DIR}/synche-cache"
    t_start=$(date +%s%N)
    manifest_id=$(run_synche "${BENCH_DIR}/source-full")
    t_end=$(date +%s%N)
    synche_ms=$(( (t_end - t_start) / 1000000 ))
    CASE1_SYNCHE+=("$synche_ms")
    echo "  >> synche:          ${synche_ms} ms"
    # Clean up: delete manifest and orphaned chunks via API.
    delete_synche_manifest "$manifest_id"

    # -- rsync (each mode) --
    for rmode in "${RSYNC_MODES[@]}"; do
        rsync_reset "$rmode"
        t_start=$(date +%s%N)
        rsync_send "$rmode" "${BENCH_DIR}/source-full" "source-full" --whole-file
        t_end=$(date +%s%N)
        rsync_ms=$(( (t_end - t_start) / 1000000 ))

        label=$(rsync_mode_label "$rmode")
        printf "  %-18s %d ms\n" "${label}:" "$rsync_ms"

        case "$rmode" in
            daemon) CASE1_RSYNC_DAEMON+=("$rsync_ms");;
            ssh)    CASE1_RSYNC_SSH+=("$rsync_ms");;
        esac
    done

    # -- scp --
    if $SCP_ENABLED; then
        scp_reset
        t_start=$(date +%s%N)
        scp_send "${BENCH_DIR}/source-full" "source-full"
        t_end=$(date +%s%N)
        scp_ms=$(( (t_end - t_start) / 1000000 ))
        CASE1_SCP+=("$scp_ms")
        echo "  >> scp:             ${scp_ms} ms"
    fi

    echo ""
done

# ---------------------------------------------------------------------------
# Case 2: Upload where half already exists on server
# ---------------------------------------------------------------------------
echo "================================================================="
echo "  CASE 2: ${SIZE_MB} MiB upload (${HALF_MB} MiB already on server)"
echo "================================================================="
echo ""

# In local mode, reset once for a clean slate before Case 2.
if [[ "$MODE" == "local" ]]; then
    reset_synche_server
fi

for ((r=1; r<=ROUNDS; r++)); do
    echo "--- Round $r/$ROUNDS ---"

    # -- synche: pre-populate server with the first half --
    rm -rf "${BENCH_DIR}/synche-cache"
    pre_manifest_id=$(run_synche "${BENCH_DIR}/preexist-half")
    # Clear local cache so the timed upload can't skip hashing.
    rm -rf "${BENCH_DIR}/synche-cache"
    # Timed: upload case-2 source (half overlaps with pre-existing chunks).
    t_start=$(date +%s%N)
    case2_manifest_id=$(run_synche "${BENCH_DIR}/source-case2")
    t_end=$(date +%s%N)
    synche_ms=$(( (t_end - t_start) / 1000000 ))
    CASE2_SYNCHE+=("$synche_ms")
    echo "  >> synche:          ${synche_ms} ms"
    # Clean up: delete both manifests (order matters — delete case2 first
    # so its unique chunks are orphaned, then delete preexist which cleans
    # up the shared chunks).
    delete_synche_manifest "$case2_manifest_id"
    delete_synche_manifest "$pre_manifest_id"

    # -- rsync (each mode): pre-populate, then delta sync --
    for rmode in "${RSYNC_MODES[@]}"; do
        rsync_reset "$rmode"
        # Pre-seed with previous version (first half matches, second half differs).
        rsync_send "$rmode" "${BENCH_DIR}/rsync-prev-case2" "source-case2" --whole-file --quiet
        # Timed: delta sync.
        t_start=$(date +%s%N)
        rsync_send "$rmode" "${BENCH_DIR}/source-case2" "source-case2"
        t_end=$(date +%s%N)
        rsync_ms=$(( (t_end - t_start) / 1000000 ))

        label=$(rsync_mode_label "$rmode")
        printf "  %-18s %d ms\n" "${label}:" "$rsync_ms"

        case "$rmode" in
            daemon) CASE2_RSYNC_DAEMON+=("$rsync_ms");;
            ssh)    CASE2_RSYNC_SSH+=("$rsync_ms");;
        esac
    done

    # -- scp (always sends full file, no delta awareness) --
    if $SCP_ENABLED; then
        scp_reset
        t_start=$(date +%s%N)
        scp_send "${BENCH_DIR}/source-case2" "source-case2"
        t_end=$(date +%s%N)
        scp_ms=$(( (t_end - t_start) / 1000000 ))
        CASE2_SCP+=("$scp_ms")
        echo "  >> scp:             ${scp_ms} ms (full file, no delta)"
    fi

    echo ""
done

# ---------------------------------------------------------------------------
# Case 3: Fresh upload of simulated disk image (100 files x 0.9 MiB)
# ---------------------------------------------------------------------------
echo "================================================================="
echo "  CASE 3: Fresh ${DISK_TOTAL_MB} MiB disk image (${DISK_FILES} x ${DISK_FILE_SIZE_KB}KB files)"
echo "================================================================="
echo ""

# In local mode, reset once for a clean slate.
if [[ "$MODE" == "local" ]]; then
    reset_synche_server
fi

for ((r=1; r<=ROUNDS; r++)); do
    echo "--- Round $r/$ROUNDS ---"

    # -- synche --
    rm -rf "${BENCH_DIR}/synche-cache"
    t_start=$(date +%s%N)
    manifest_id=$(run_synche "${BENCH_DIR}/disk-original")
    t_end=$(date +%s%N)
    synche_ms=$(( (t_end - t_start) / 1000000 ))
    CASE3_SYNCHE+=("$synche_ms")
    echo "  >> synche:          ${synche_ms} ms"
    delete_synche_manifest "$manifest_id"

    # -- rsync (each mode) --
    for rmode in "${RSYNC_MODES[@]}"; do
        rsync_reset "$rmode"
        t_start=$(date +%s%N)
        rsync_send "$rmode" "${BENCH_DIR}/disk-original" "disk-original" --whole-file
        t_end=$(date +%s%N)
        rsync_ms=$(( (t_end - t_start) / 1000000 ))

        label=$(rsync_mode_label "$rmode")
        printf "  %-18s %d ms\n" "${label}:" "$rsync_ms"

        case "$rmode" in
            daemon) CASE3_RSYNC_DAEMON+=("$rsync_ms");;
            ssh)    CASE3_RSYNC_SSH+=("$rsync_ms");;
        esac
    done

    # -- scp --
    if $SCP_ENABLED; then
        scp_reset
        t_start=$(date +%s%N)
        scp_send "${BENCH_DIR}/disk-original" "disk-original"
        t_end=$(date +%s%N)
        scp_ms=$(( (t_end - t_start) / 1000000 ))
        CASE3_SCP+=("$scp_ms")
        echo "  >> scp:             ${scp_ms} ms"
    fi

    echo ""
done

# ---------------------------------------------------------------------------
# Case 4: Re-upload disk image after 10% of files changed
# ---------------------------------------------------------------------------
echo "================================================================="
echo "  CASE 4: ${DISK_TOTAL_MB} MiB disk image, ${DISK_CHANGED}/${DISK_FILES} files changed (~${DISK_CHANGED_KB}KB delta)"
echo "================================================================="
echo ""

# In local mode, reset once for a clean slate.
if [[ "$MODE" == "local" ]]; then
    reset_synche_server
fi

for ((r=1; r<=ROUNDS; r++)); do
    echo "--- Round $r/$ROUNDS ---"

    # -- synche: pre-populate with original disk image --
    rm -rf "${BENCH_DIR}/synche-cache"
    orig_manifest_id=$(run_synche "${BENCH_DIR}/disk-original")
    rm -rf "${BENCH_DIR}/synche-cache"
    # Timed: upload modified disk image (90% chunks already exist).
    t_start=$(date +%s%N)
    mod_manifest_id=$(run_synche "${BENCH_DIR}/disk-modified")
    t_end=$(date +%s%N)
    synche_ms=$(( (t_end - t_start) / 1000000 ))
    CASE4_SYNCHE+=("$synche_ms")
    echo "  >> synche:          ${synche_ms} ms"
    delete_synche_manifest "$mod_manifest_id"
    delete_synche_manifest "$orig_manifest_id"

    # -- rsync (each mode): pre-seed with original, then delta-sync modified --
    for rmode in "${RSYNC_MODES[@]}"; do
        rsync_reset "$rmode"
        rsync_send "$rmode" "${BENCH_DIR}/rsync-prev-case4" "disk-image" --whole-file --quiet
        t_start=$(date +%s%N)
        rsync_send "$rmode" "${BENCH_DIR}/disk-modified" "disk-image"
        t_end=$(date +%s%N)
        rsync_ms=$(( (t_end - t_start) / 1000000 ))

        label=$(rsync_mode_label "$rmode")
        printf "  %-18s %d ms\n" "${label}:" "$rsync_ms"

        case "$rmode" in
            daemon) CASE4_RSYNC_DAEMON+=("$rsync_ms");;
            ssh)    CASE4_RSYNC_SSH+=("$rsync_ms");;
        esac
    done

    # -- scp (always sends full file) --
    if $SCP_ENABLED; then
        scp_reset
        t_start=$(date +%s%N)
        scp_send "${BENCH_DIR}/disk-modified" "disk-modified"
        t_end=$(date +%s%N)
        scp_ms=$(( (t_end - t_start) / 1000000 ))
        CASE4_SCP+=("$scp_ms")
        echo "  >> scp:             ${scp_ms} ms (full file, no delta)"
    fi

    echo ""
done

# ---------------------------------------------------------------------------
# Case 5: Fresh directory upload (100 individual files)
# ---------------------------------------------------------------------------
echo "================================================================="
echo "  CASE 5: Fresh directory upload (${DISK_FILES} x ${DISK_FILE_SIZE_KB}KB files, ~${DISK_TOTAL_MB} MiB)"
echo "================================================================="
echo ""

# In local mode, reset once for a clean slate.
if [[ "$MODE" == "local" ]]; then
    reset_synche_server
fi

for ((r=1; r<=ROUNDS; r++)); do
    echo "--- Round $r/$ROUNDS ---"

    # -- synche --
    rm -rf "${BENCH_DIR}/synche-cache"
    t_start=$(date +%s%N)
    manifest_ids=$(run_synche_dir "${BENCH_DIR}/dir-original")
    t_end=$(date +%s%N)
    synche_ms=$(( (t_end - t_start) / 1000000 ))
    CASE5_SYNCHE+=("$synche_ms")
    echo "  >> synche:          ${synche_ms} ms"
    delete_synche_manifests "$manifest_ids"

    # -- rsync (each mode) --
    for rmode in "${RSYNC_MODES[@]}"; do
        rsync_reset "$rmode"
        t_start=$(date +%s%N)
        rsync_send_dir "$rmode" "${BENCH_DIR}/dir-original" --whole-file
        t_end=$(date +%s%N)
        rsync_ms=$(( (t_end - t_start) / 1000000 ))

        label=$(rsync_mode_label "$rmode")
        printf "  %-18s %d ms\n" "${label}:" "$rsync_ms"

        case "$rmode" in
            daemon) CASE5_RSYNC_DAEMON+=("$rsync_ms");;
            ssh)    CASE5_RSYNC_SSH+=("$rsync_ms");;
        esac
    done

    # -- scp --
    if $SCP_ENABLED; then
        scp_reset
        t_start=$(date +%s%N)
        scp_send_dir "${BENCH_DIR}/dir-original"
        t_end=$(date +%s%N)
        scp_ms=$(( (t_end - t_start) / 1000000 ))
        CASE5_SCP+=("$scp_ms")
        echo "  >> scp:             ${scp_ms} ms"
    fi

    echo ""
done

# ---------------------------------------------------------------------------
# Case 6: Re-upload directory after 10% of files changed
# ---------------------------------------------------------------------------
echo "================================================================="
echo "  CASE 6: Directory re-upload, ${DISK_CHANGED}/${DISK_FILES} files changed (~${DISK_CHANGED_KB}KB delta)"
echo "================================================================="
echo ""

# In local mode, reset once for a clean slate.
if [[ "$MODE" == "local" ]]; then
    reset_synche_server
fi

for ((r=1; r<=ROUNDS; r++)); do
    echo "--- Round $r/$ROUNDS ---"

    # -- synche: pre-populate with original directory --
    rm -rf "${BENCH_DIR}/synche-cache"
    orig_manifest_ids=$(run_synche_dir "${BENCH_DIR}/dir-original")
    rm -rf "${BENCH_DIR}/synche-cache"
    # Timed: upload modified directory (90% of files unchanged, chunks already on server).
    t_start=$(date +%s%N)
    mod_manifest_ids=$(run_synche_dir "${BENCH_DIR}/dir-modified")
    t_end=$(date +%s%N)
    synche_ms=$(( (t_end - t_start) / 1000000 ))
    CASE6_SYNCHE+=("$synche_ms")
    echo "  >> synche:          ${synche_ms} ms"
    delete_synche_manifests "$mod_manifest_ids"
    delete_synche_manifests "$orig_manifest_ids"

    # -- rsync (each mode): pre-seed with original, then sync modified --
    for rmode in "${RSYNC_MODES[@]}"; do
        rsync_reset "$rmode"
        rsync_send_dir "$rmode" "${BENCH_DIR}/dir-original" --whole-file --quiet
        t_start=$(date +%s%N)
        rsync_send_dir "$rmode" "${BENCH_DIR}/dir-modified"
        t_end=$(date +%s%N)
        rsync_ms=$(( (t_end - t_start) / 1000000 ))

        label=$(rsync_mode_label "$rmode")
        printf "  %-18s %d ms\n" "${label}:" "$rsync_ms"

        case "$rmode" in
            daemon) CASE6_RSYNC_DAEMON+=("$rsync_ms");;
            ssh)    CASE6_RSYNC_SSH+=("$rsync_ms");;
        esac
    done

    # -- scp (always sends everything) --
    if $SCP_ENABLED; then
        scp_reset
        t_start=$(date +%s%N)
        scp_send_dir "${BENCH_DIR}/dir-modified"
        t_end=$(date +%s%N)
        scp_ms=$(( (t_end - t_start) / 1000000 ))
        CASE6_SCP+=("$scp_ms")
        echo "  >> scp:             ${scp_ms} ms (full dir, no delta)"
    fi

    echo ""
done

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
echo "================================================================="
echo "  RESULTS SUMMARY"
echo "================================================================="
echo ""
if [[ "$MODE" == "remote" ]]; then
    echo "  server: ${SERVER_HOST}"
    echo ""
fi

avg() {
    local sum=0 count=0
    for v in "$@"; do
        sum=$((sum + v))
        count=$((count + 1))
    done
    echo $((sum / count))
}

c1_synche_avg=$(avg "${CASE1_SYNCHE[@]}")
c2_synche_avg=$(avg "${CASE2_SYNCHE[@]}")
c3_synche_avg=$(avg "${CASE3_SYNCHE[@]}")
c4_synche_avg=$(avg "${CASE4_SYNCHE[@]}")
c5_synche_avg=$(avg "${CASE5_SYNCHE[@]}")
c6_synche_avg=$(avg "${CASE6_SYNCHE[@]}")

# Build competitor columns dynamically: rsync modes + scp.
declare -a all_labels c1_all_avgs c2_all_avgs c3_all_avgs c4_all_avgs c5_all_avgs c6_all_avgs

for rmode in "${RSYNC_MODES[@]}"; do
    label=$(rsync_mode_label "$rmode")
    all_labels+=("$label")

    case "$rmode" in
        daemon)
            c1_all_avgs+=($(avg "${CASE1_RSYNC_DAEMON[@]}"))
            c2_all_avgs+=($(avg "${CASE2_RSYNC_DAEMON[@]}"))
            c3_all_avgs+=($(avg "${CASE3_RSYNC_DAEMON[@]}"))
            c4_all_avgs+=($(avg "${CASE4_RSYNC_DAEMON[@]}"))
            c5_all_avgs+=($(avg "${CASE5_RSYNC_DAEMON[@]}"))
            c6_all_avgs+=($(avg "${CASE6_RSYNC_DAEMON[@]}"))
            ;;
        ssh)
            c1_all_avgs+=($(avg "${CASE1_RSYNC_SSH[@]}"))
            c2_all_avgs+=($(avg "${CASE2_RSYNC_SSH[@]}"))
            c3_all_avgs+=($(avg "${CASE3_RSYNC_SSH[@]}"))
            c4_all_avgs+=($(avg "${CASE4_RSYNC_SSH[@]}"))
            c5_all_avgs+=($(avg "${CASE5_RSYNC_SSH[@]}"))
            c6_all_avgs+=($(avg "${CASE6_RSYNC_SSH[@]}"))
            ;;
    esac
done

if $SCP_ENABLED && [[ ${#CASE1_SCP[@]} -gt 0 ]]; then
    all_labels+=("scp")
    c1_all_avgs+=($(avg "${CASE1_SCP[@]}"))
    c2_all_avgs+=($(avg "${CASE2_SCP[@]}"))
    c3_all_avgs+=($(avg "${CASE3_SCP[@]}"))
    c4_all_avgs+=($(avg "${CASE4_SCP[@]}"))
    c5_all_avgs+=($(avg "${CASE5_SCP[@]}"))
    c6_all_avgs+=($(avg "${CASE6_SCP[@]}"))
fi

# Print table header.
printf "  %-50s %10s" "" "synche"
for label in "${all_labels[@]}"; do
    printf "  %14s" "$label"
done
echo ""

printf "  %-50s %10s" "" "--------"
for _ in "${all_labels[@]}"; do
    printf "  %14s" "--------------"
done
echo ""

printf "  %-50s %7d ms" "Case 1: ${SIZE_MB} MiB fresh upload" "$c1_synche_avg"
for val in "${c1_all_avgs[@]}"; do printf "  %11d ms" "$val"; done
echo ""

printf "  %-50s %7d ms" "Case 2: ${SIZE_MB} MiB, ${HALF_MB} MiB exists" "$c2_synche_avg"
for val in "${c2_all_avgs[@]}"; do printf "  %11d ms" "$val"; done
echo ""

printf "  %-50s %7d ms" "Case 3: ${DISK_TOTAL_MB} MiB disk image fresh" "$c3_synche_avg"
for val in "${c3_all_avgs[@]}"; do printf "  %11d ms" "$val"; done
echo ""

printf "  %-50s %7d ms" "Case 4: disk image, 10%% changed" "$c4_synche_avg"
for val in "${c4_all_avgs[@]}"; do printf "  %11d ms" "$val"; done
echo ""

printf "  %-50s %7d ms" "Case 5: directory (${DISK_FILES} files) fresh" "$c5_synche_avg"
for val in "${c5_all_avgs[@]}"; do printf "  %11d ms" "$val"; done
echo ""

printf "  %-50s %7d ms" "Case 6: directory, 10%% files changed" "$c6_synche_avg"
for val in "${c6_all_avgs[@]}"; do printf "  %11d ms" "$val"; done
echo ""
echo ""

# Speed comparisons — synche vs each competitor.
declare -a case_labels case_synche_avgs
case_labels=("Case 1" "Case 2" "Case 3" "Case 4" "Case 5" "Case 6")
case_synche_avgs=("$c1_synche_avg" "$c2_synche_avg" "$c3_synche_avg" "$c4_synche_avg" "$c5_synche_avg" "$c6_synche_avg")

for i in "${!all_labels[@]}"; do
    label="${all_labels[$i]}"
    declare -a comp_avgs=("${c1_all_avgs[$i]}" "${c2_all_avgs[$i]}" "${c3_all_avgs[$i]}" "${c4_all_avgs[$i]}" "${c5_all_avgs[$i]}" "${c6_all_avgs[$i]}")
    for c in 0 1 2 3 4 5; do
        s_avg=${case_synche_avgs[$c]}
        r_avg=${comp_avgs[$c]}
        clabel=${case_labels[$c]}
        if [[ "$s_avg" -lt "$r_avg" && "$s_avg" -gt 0 ]]; then
            speedup=$(echo "scale=1; $r_avg / $s_avg" | bc)
            echo "  ${clabel} vs ${label}: synche is ${speedup}x faster"
        elif [[ "$r_avg" -lt "$s_avg" && "$r_avg" -gt 0 ]]; then
            speedup=$(echo "scale=1; $s_avg / $r_avg" | bc)
            echo "  ${clabel} vs ${label}: ${label} is ${speedup}x faster"
        else
            echo "  ${clabel} vs ${label}: roughly tied"
        fi
    done
done
echo ""

# Individual rounds — helper function.
print_rounds() {
    local case_label="$1"
    shift
    local -n synche_arr=$1; shift
    local -n daemon_arr=$1; shift
    local -n ssh_arr=$1; shift
    local -n scp_arr=$1; shift

    echo "  ${case_label} individual rounds (ms):"
    for ((i=0; i<ROUNDS; i++)); do
        printf "    round %d:  synche=%d" $((i+1)) "${synche_arr[$i]}"
        for rmode in "${RSYNC_MODES[@]}"; do
            case "$rmode" in
                daemon) printf "  daemon=%d" "${daemon_arr[$i]}";;
                ssh)    printf "  ssh=%d" "${ssh_arr[$i]}";;
            esac
        done
        if $SCP_ENABLED && [[ ${#scp_arr[@]} -gt 0 ]]; then
            printf "  scp=%d" "${scp_arr[$i]}"
        fi
        echo ""
    done
    echo ""
}

print_rounds "Case 1" CASE1_SYNCHE CASE1_RSYNC_DAEMON CASE1_RSYNC_SSH CASE1_SCP
print_rounds "Case 2" CASE2_SYNCHE CASE2_RSYNC_DAEMON CASE2_RSYNC_SSH CASE2_SCP
print_rounds "Case 3" CASE3_SYNCHE CASE3_RSYNC_DAEMON CASE3_RSYNC_SSH CASE3_SCP
print_rounds "Case 4" CASE4_SYNCHE CASE4_RSYNC_DAEMON CASE4_RSYNC_SSH CASE4_SCP
print_rounds "Case 5" CASE5_SYNCHE CASE5_RSYNC_DAEMON CASE5_RSYNC_SSH CASE5_SCP
print_rounds "Case 6" CASE6_SYNCHE CASE6_RSYNC_DAEMON CASE6_RSYNC_SSH CASE6_SCP

echo "================================================================="
