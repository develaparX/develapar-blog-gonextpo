import { useParams } from "react-router-dom";

const ArticleDetailPage = () => {
  const { slug } = useParams();

  return (
    <div className="max-w-3xl mx-auto px-4 py-12">
      <article>
        <h1 className="text-4xl font-bold mb-4 leading-tight">
          Judul Artikel: {slug?.replace(/-/g, " ")}
        </h1>

        <div className="text-gray-500 text-sm mb-8">
          By <span className="font-semibold">John Doe</span> · June 13, 2025 · 5
          min read
        </div>

        <img
          src="https://source.unsplash.com/random/900x400?article"
          alt="Article Banner"
          className="w-full h-72 object-cover rounded-lg mb-8"
        />

        <div className="prose prose-lg max-w-none">
          <p>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer nec
            odio. Praesent libero. Sed cursus ante dapibus diam. Sed nisi. Nulla
            quis sem at nibh elementum imperdiet. Duis sagittis ipsum.
          </p>
          <p>
            Praesent mauris. Fusce nec tellus sed augue semper porta. Mauris
            massa. Vestibulum lacinia arcu eget nulla. Class aptent taciti
            sociosqu ad litora torquent per conubia nostra, per inceptos
            himenaeos.
          </p>
          <h2>Subjudul dalam Artikel</h2>
          <p>
            Curabitur sodales ligula in libero. Sed dignissim lacinia nunc.
            Curabitur tortor. Pellentesque nibh. Aenean quam. In scelerisque sem
            at dolor. Maecenas mattis.
          </p>
          <blockquote>
            "Artikel yang baik harus memancing rasa penasaran dan menjawabnya
            secara elegan."
          </blockquote>
          <p>
            Sed convallis tristique sem. Proin ut ligula vel nunc egestas
            porttitor. Morbi lectus risus, iaculis vel, suscipit quis, luctus
            non, massa. Fusce ac turpis quis ligula lacinia aliquet.
          </p>
        </div>
      </article>
    </div>
  );
};

export default ArticleDetailPage;
