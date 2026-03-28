import "./globals.css";

export const metadata = {
  title: "Project Dashboard",
  description: "Mini project management dashboard",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
