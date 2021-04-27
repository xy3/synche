import Link from "next/link";

interface BreadcrumbLink {
  href?: string;
  title: string;
}

interface ComponentProps {
  links: Array<BreadcrumbLink>;
}

export default function Breadcrumb({ links }: ComponentProps) {
  return (
    <section className="w-full flex my-4">
      {links.map((link: BreadcrumbLink, idx) => {
        return (
          <div className="flex items-center flex-wrap" key={idx}>
            {link.href ? (
              <>
                <Link key={idx} href={link.href}>
                  <a className="inline mx-2 text-gray-500 truncate text-sm">
                    {link.title}
                  </a>
                </Link>
                {idx !== links.length - 1 ? (
                  <span className="select-none">&gt;</span>
                ) : null}
              </>
            ) : (
              <span className="text-sm inline mx-2 text-gray-400">
                {link.title}
              </span>
            )}
          </div>
        );
      })}
    </section>
  );
}
