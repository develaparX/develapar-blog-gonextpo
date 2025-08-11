import {
  Card,
  CardAction,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ArrowRight } from "lucide-react";

const dummyData = [
  {
    title: "lorem ipsum dolor sit amet consectetur adipisicing elit.",
    description:
      "lorem ipsum dolor sit amet consectetur adipisicing elit. Iusto corrupti impedit ipsam,",
    author: "Author 1",
    imageUrl:
      "https://images.unsplash.com/photo-1743449661678-c22cd73b338a?q=80&w=870&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
  },
  {
    title: "lorem ipsum dolor sit amet consectetur adipisicing elit.",
    description:
      "lorem ipsum dolor sit amet consectetur adipisicing elit. Iusto corrupti impedit ipsam,",
    author: "Author 2",

    imageUrl:
      "https://images.unsplash.com/photo-1743449661678-c22cd73b338a?q=80&w=870&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
  },
  {
    title: "lorem ipsum dolor sit amet consectetur adipisicing elit.",
    description:
      "lorem ipsum dolor sit amet consectetur adipisicing elit. Iusto corrupti impedit ipsam,",
    author: "Author 3",

    imageUrl:
      "https://images.unsplash.com/photo-1743449661678-c22cd73b338a?q=80&w=870&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
  },
];

const Homepage = () => {
  return (
    <div className="mx-auto px-5">
      <div className="my-10 flex justify-center items-center w-full text-4xl font-extrabold">
        OH Hi! Welcome to develapar!
      </div>

      {/* TOP HERO */}
      <div
        id="top-menu"
        className="flex flex-row gap-5 justify-between w-full h-[40vw] items-center "
      >
        {/* === KARTU KIRI YANG DIMODIFIKASI UNTUK REACT === */}
        <Card
          className="group relative w-full h-full max-w-[70%] rounded-4xl p-0 shadow-none overflow-hidden
                     transition-all duration-300 ease-in-out "
        >
          <img
            alt="Background"
            src="https://images.unsplash.com/photo-1743449661678-c22cd73b338a?q=80&w=870&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D" // Ganti
            className="absolute inset-0 w-full h-full object-cover"
          />

          <div
            className="absolute inset-0 bg-black 
                       opacity-50 group-hover:opacity-0 
                       transition-opacity duration-300 ease-in-out"
          />

          <div className="relative flex flex-col h-full justify-end">
            <CardHeader className="flex flex-row items-end justify-between p-0 text-white">
              <div className="flex flex-col gap-1 p-5">
                <CardTitle className="text-2xl font-bold">
                  Lorem ipsum dolor, sit amet consectetur adipisicing elit.
                </CardTitle>
                <CardDescription className="text-gray-300">
                  Ini Decsripsi conetenya apa Lorem ipsum dolor sit, amet
                  consectetur adipisicing elit. Iusto corrupti impedit ipsam,
                </CardDescription>
                <CardDescription className="text-gray-300">
                  Authornya disini bro
                </CardDescription>
              </div>

              <CardAction className="bg-white text-black min-w-[20%] h-auto flex items-center justify-center rounded-lg shadow-none m-5 p-4 hover:scale-110">
                <ArrowRight size={80} />
              </CardAction>
            </CardHeader>
          </div>
        </Card>

        {/* Kanan */}
        <div className="w-full h-full max-w-[35%] flex flex-col gap-5 justify-between">
          {dummyData.map((item, index) => (
            <Card
              key={index}
              className="py-2 shadow-none border-none flex flex-row h-1/3"
            >
              <Card className="relative group w-[50%] h-[100%] rounded-xl overflow-hidden">
                <img
                  src={item.imageUrl}
                  alt={item.title}
                  className="absolute inset-0 w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
                />

                <div
                  className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent 
               opacity-70 group-hover:opacity-100 transition-opacity duration-300"
                />

                <CardAction
                  className="absolute bottom-4 right-4 bg-white text-black rounded-full p-2 m-0
              opacity-100 transition-all duration-300 ease-in-out
               transform group-hover:scale-110"
                >
                  <ArrowRight size={24} />
                </CardAction>
              </Card>

              <CardHeader className="flex px-0 flex-col justify-evenly items-start w-full">
                <CardTitle className="pt-1">{item.title}</CardTitle>
                <CardDescription>{item.description}</CardDescription>
                <CardDescription className="pb-1">
                  By: {item.author}
                </CardDescription>
              </CardHeader>
            </Card>
          ))}
        </div>
      </div>

      {/* MIDDLE HERO */}
    </div>
  );
};

export default Homepage;
