import jwt_decode from "jwt-decode";

interface IDecodedUser {
  UserID: number;
  Email: string;
  Name: string;
  Picture: string;
  Role: string;
  TokenType: string;
}

function isLoggedIn(token: string): boolean {
  if (!token) {
    return false;
  }

  try {
    const decodedToken: object = jwt_decode(token);

    return (
      decodedToken &&
      decodedToken.hasOwnProperty("UserID") &&
      decodedToken.hasOwnProperty("Email") &&
      decodedToken.hasOwnProperty("Name") &&
      decodedToken.hasOwnProperty("Picture") &&
      decodedToken.hasOwnProperty("Role") &&
      decodedToken.hasOwnProperty("TokenType")
    );
  } catch (err) {
    return false;
  }
}

function decodeToken(token: string): IDecodedUser {
  try {
    const decodedToken: IDecodedUser = jwt_decode(token);
    return decodedToken;
  } catch (err) {
    return null;
  }
}

export { isLoggedIn, decodeToken };
