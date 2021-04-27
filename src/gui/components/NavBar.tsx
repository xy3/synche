import Link from "next/link";
import { UserConsumer } from "../context/userContext";

export default function NavBar() {
  return (
    <nav className="w-full flex justify-between items-center p-4 bg-white shadow-sm">
      <div>
        <Link href="/">
          <a>
            <img src="/img/logo.png" className="w-32 h-auto" />
          </a>
        </Link>
      </div>
      <UserConsumer>
        {(user) =>
          user.loggedIn ? (
            <div className="flex">
              <Link href="/dashboard">
                <a className="font-bold mx-2 p-2 text-brand-dark-blue border-b border-transparent hover:border-brand-blue hover:text-brand-blue">
                  Dashboard
                </a>
              </Link>
              <Link href="/logout">
                <a className="mx-2 p-2 text-red-600 border-b border-transparent hover:border-red-500 hover:text-red-500">
                  Log Out
                </a>
              </Link>
            </div>
          ) : (
            <div className="flex">
              <Link href="/login">
                <a className="mx-2 p-2 text-brand-dark-blue border-b border-transparent hover:border-brand-blue hover:text-brand-blue">
                  Log In
                </a>
              </Link>
              <Link href="/signup">
                <a className="font-bold mx-2 p-2 text-brand-blue border border-brand-blue">
                  Sign Up
                </a>
              </Link>
            </div>
          )
        }
      </UserConsumer>
    </nav>
  );
}
