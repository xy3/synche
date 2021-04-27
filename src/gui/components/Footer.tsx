import Link from "next/link";

export default function Footer() {
  return (
    <footer className="mt-auto w-full py-16 bg-brand-dark-blue">
      <div className="container p-4 flex justify-between items-center">
        <img src="/img/logo-white.png" className="w-24 md:w-48 h-auto" />
        <div className="flex items-center">
          <Link href="/dashboard">
            <a className="text-gray-400 mx-2">Dashboard</a>
          </Link>
          <a className="text-gray-400 mx-2" href="#">
            API
          </a>
          <a className="text-gray-400 mx-2" href="/#contact">
            Contact Us
          </a>
        </div>
      </div>
    </footer>
  );
}
