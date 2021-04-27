import Head from "next/head";
import { ReactChild } from "react";
import Footer from "./Footer";
import NavBar from "./NavBar";

interface ComponentProps {
  title?: string;
  children: ReactChild | ReactChild[];
}

export default function Layout({ title, children }: ComponentProps) {
  return (
    <>
      <Head>
        <link
          rel="icon"
          type="image/png"
          href="favicon-32x32.png"
          sizes="32x32"
        />
        <link
          rel="icon"
          type="image/png"
          href="favicon-16x16.png"
          sizes="16x16"
        />

        <link rel="preconnect" href="https://fonts.gstatic.com" />
        <link
          href="https://fonts.googleapis.com/css2?family=Rubik:wght@400;700&family=Noto+Sans+KR:wght@400;700&display=swap"
          rel="stylesheet"
        />
        <title>Synche{title ? `: ${title}` : null}</title>
      </Head>
      <div className="flex flex-col" style={{ height: "100vh", margin: 0 }}>
        <NavBar />
        <div className="mb-8">{children}</div>
        <Footer />
      </div>
    </>
  );
}
