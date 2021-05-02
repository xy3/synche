import { createContext, useEffect, useState } from "react";
import Cookies from "js-cookie";
import jwt_decode from "jwt-decode";

const UserContext = createContext(null);

interface IUser {
  UserID?: number;
  Email?: string;
  Name?: string;
  Picture?: string;
  Role?: string;
  TokenType?: string;
  loggedIn: boolean;
}

function UserProvider({ children }) {
  const [userState, setUserState] = useState<IUser>({
    loggedIn: false,
  });

  useEffect(() => {
    async function loadUser() {
      try {
        if (Cookies.get("accessToken")) {
          const decodedToken: object = jwt_decode(Cookies.get("accessToken"));

          if (
            decodedToken.hasOwnProperty("UserID") &&
            decodedToken.hasOwnProperty("Email") &&
            decodedToken.hasOwnProperty("Name") &&
            decodedToken.hasOwnProperty("Picture") &&
            decodedToken.hasOwnProperty("Role") &&
            decodedToken.hasOwnProperty("TokenType")
          ) {
            setUserState({
              ...decodedToken,
              loggedIn: true,
            });
          } else {
            Cookies.remove("token");
          }
        }
      } catch (err) {}
    }

    loadUser();
  }, []);

  return (
    <UserContext.Provider value={userState}>{children}</UserContext.Provider>
  );
}

const UserConsumer = UserContext.Consumer;

export default UserProvider;
export { UserConsumer, UserContext };
