export default function Page({ params }: { params: { id: string } }) {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <h1>Order {params.id}</h1>
    </main>
  );
}
