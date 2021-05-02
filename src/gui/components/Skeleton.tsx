interface ComponentProps {
  message: string;
}

export default function Skeleton({ message }: ComponentProps) {
  return (
    <div className="my-4 p-4 rounded-sm border border-gray-400">
      <p className="text-center text-gray-400">{message}</p>
    </div>
  );
}
