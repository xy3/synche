import jwt_decode from "jwt-decode";

function isLoggedIn(token: string): boolean {
  if (!token) {
    return false;
  }

  return true;
  /*try {
    const decodedToken: any = jwt_decode(token);

    return (
      decodedToken
    );
  } catch (err) {
    return false;
  }*/
}

export { isLoggedIn };
