import {
  Card,
  CardAction,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ArrowLeft, ArrowRight, ArrowUpRight } from "lucide-react";
import { useRef } from "react";

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

const dummyNewsData = [
  {
    title:
      "Rare Solar Flare Sparks Global Auroras and Disrupts Satellite Signals",
    description:
      "A solar flare caused rare auroras to appear much farther south than usual, delighting skywatchers.",
    authorName: "Lina Carter",
    authorAvatar: "/path/to/avatar1.jpg", // Ganti dengan path avatar
    imageUrl:
      "https://images.unsplash.com/photo-1558981403-c5f9899a28bc?q=80&w=2070&auto=format&fit=crop",
  },
  {
    title: "AI Tutors Launch in Rural Schools, Transforming Education Access",
    description:
      "A tech company has launched AI-powered tutors that deliver personalized lessons to students in rural schools.",
    authorName: "Rahul Desai",
    authorAvatar: "/path/to/avatar2.jpg",
    imageUrl:
      "https://images.unsplash.com/photo-1558981403-c5f9899a28bc?q=80&w=2070&auto=format&fit=crop",
  },
  {
    title: "Historic Climate Pact Reached After 18-Hour Negotiation Marathon",
    description:
      "Leaders from over 40 nations have signed a new climate accord aimed at cutting global carbon emissions.",
    authorName: "Serena Ngo",
    authorAvatar: "/path/to/avatar3.jpg",
    imageUrl:
      "https://images.unsplash.com/photo-1558981403-c5f9899a28bc?q=80&w=2070&auto=format&fit=crop",
  },
  {
    title: "Mysterious Album Teased with Midnight Countdown",
    description:
      "A surprise album announcement has taken the music world by storm after a major label dropped a cryptic timer.",
    authorName: "J. Lane",
    authorAvatar: "/path/to/avatar4.jpg",
    imageUrl:
      "https://images.unsplash.com/photo-1516450360452-9312f5e86fc7?q=80&w=2070&auto=format&fit=crop",
  },
];

const mainPopularStory = {
  title:
    "Construction Begins on International Space Agencies' New Orbital Station",
  description:
    "A collaborative initiative among major international space agencies has officially commenced the construction of a next-generation orbital space station, designed to serve as a critical hub for future deep space exploration...",
  authorName: "Marcus Iliev",
  authorAvatar:
    "https://images.unsplash.com/photo-1558507652-2d9626c4e67a?q=80&w=1974&auto=format&fit=crop",
  imageUrl:
    "https://images.unsplash.com/photo-1558981403-c5f9899a28bc?q=80&w=2070&auto=format&fit=crop",
};

const secondaryPopularStories = [
  {
    title: "Ten Smart Cities Go Fully Digital to Improve Daily Urban Life",
    authorName: "James Okoro",
    authorAvatar: "/path/to/avatar6.jpg",
    imageUrl:
      "https://images.unsplash.com/photo-1558507652-2d9626c4e67a?q=80&w=1974&auto=format&fit=crop",
  },
  {
    title: "Vertical Farms Bring Fresh Produce to Remote Regions",
    authorName: "Leo Mmatha",
    authorAvatar: "/path/to/avatar7.jpg",
    imageUrl:
      "https://images.unsplash.com/photo-1558507652-2d9626c4e67a?q=80&w=1974&auto=format&fit=crop",
  },
  {
    title: "Virtual Reality History Lessons Engage Students in New Ways",
    authorName: "Henrik Sorensen",
    authorAvatar:
      "https://images.unsplash.com/photo-1558507652-2d9626c4e67a?q=80&w=1974&auto=format&fit=crop",
    imageUrl:
      "https://images.unsplash.com/photo-1558507652-2d9626c4e67a?q=80&w=1974&auto=format&fit=crop",
  },
];

// Data dummy untuk konten highlight
const mainHighlightStory = {
  title:
    "Smart Clothing Tracks Health Stats and Sends Alerts to Users in Real Time",
  description:
    "New wearable technology, seamlessly woven into everyday clothing, is revolutionizing personal health monitoring. This advanced fabric technology tracks key metrics such as heart rate, stress levels, hydration, and sleep cycles in real time.",
  authorName: "J. Lane, Staff Writer",
  authorAvatar: "/path/to/avatar9.jpg", // Ganti dengan path avatar
  imageUrl:
    "https://images.unsplash.com/photo-1558981403-c5f9899a28bc?q=80&w=2070&auto=format&fit=crop",
};

const secondaryHighlightStories = [
  {
    title: "Global Solar Power Output Hits New Record, Surpassing 1 Terawatt",
    authorName: "Emma Torres",
    authorAvatar: "/path/to/avatar10.jpg",
    imageUrl:
      "https://images.unsplash.com/photo-1508515329313-a44d189bb396?q=80&w=2071&auto=format&fit=crop",
  },
  {
    title:
      "Self-Healing Concrete Developed to Extend the Life of Infrastructure Projects",
    authorName: "Greg Porter",
    authorAvatar: "/path/to/avatar11.jpg",
    imageUrl:
      "https://images.unsplash.com/photo-1581093192089-14a4a083a3a4?q=80&w=2070&auto=format&fit=crop",
  },
  {
    title:
      "Smart Hospitals Use AI to Predict Patient Needs and Prevent Emergencies",
    authorName: "Tanya Patel",
    authorAvatar: "/path/to/avatar12.jpg",
    imageUrl:
      "https://images.unsplash.com/photo-1538108149393-fbbd81895907?q=80&w=2128&auto=format&fit=crop",
  },
  {
    title: "Electric Cars Outnumber Gas-Powered Vehicles in Major City Center",
    authorName: "Luis Garza",
    authorAvatar: "/path/to/avatar13.jpg",
    imageUrl:
      "https://images.unsplash.com/photo-1534294643621-e9b4164b3554?q=80&w=2070&auto=format&fit=crop",
  },
];

const Homepage = () => {
  const scrollContainerRef = useRef<HTMLDivElement>(null);

  // 3. Buat fungsi untuk menggeser kontainer
  const handleScroll = (scrollOffset: any) => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollBy({
        left: scrollOffset,
        behavior: "smooth", // Efek scroll yang mulus
      });
    }
  };
  return (
    <div className="mx-auto">
      <div className="my-10 flex justify-center items-center w-full text-4xl font-extrabold">
        OH Hi! Welcome to develapar!
      </div>

      {/* TOP HERO */}
      <div
        id="top-menu"
        className="flex flex-row gap-5 justify-between w-full h-[40vw] items-center px-5"
      >
        {/* === KARTU KIRI YANG DIMODIFIKASI UNTUK REACT === */}
        <Card className="group relative w-full h-full max-w-[60%]  p-8 overflow-visible transition-all duration-300 ease-in-out">
          <img
            alt="Background"
            src="https://images.unsplash.com/photo-1743449661678-c22cd73b338a?q=80&w=870&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
            className="absolute inset-0 w-full h-full object-cover rounded-xl "
          />

          <div
            className="absolute inset-0 bg-gradient-to-t from-black/100 to-transparent
               opacity-100 group-hover:opacity-100 transition-opacity duration-300 rounded-xl"
          />

          <div className="relative flex flex-col h-full justify-end text-white max-w-[75%]">
            <h2 className="text-2xl font-bold">
              Lorem ipsum dolor, sit amet consectetur adipisicing elit.
            </h2>
            <p className="mt-2 text-gray-300">
              Ini Decsripsi conetenya apa Lorem ipsum dolor sit, amet
              consectetur adipisicing elit. Iusto corrupti impedit ipsam,
            </p>
            <p className="text-gray-300">Authornya disini bro</p>
          </div>
          <CardAction className="group absolute -bottom-2 -right-3 bg-white p-5 rounded-tl-[40px] rounded-tr-2xl rounded-bl-2xl  transition-transform duration-300">
            <div className="flex p-4 items-center justify-center w-full h-full  border-8 border-blue-400 bg-blue-800 rounded-4xl transition-transform duration-300 group-hover:scale-110">
              <ArrowRight size={64} className=" text-white" />
            </div>
          </CardAction>
        </Card>

        {/* Kanan */}
        <div className="w-full h-full max-w-[40%] flex flex-col gap-5 justify-between ">
          {dummyData.map((item, index) => (
            <Card
              key={index}
              className="py-2 shadow-none border-none flex rounded-xl flex-row h-1/3"
            >
              <Card className="relative group w-[50%] h-[100%] rounded-xl overflow-visible">
                <img
                  src={item.imageUrl}
                  alt={item.title}
                  className="absolute inset-0 w-full h-full object-cover rounded-xl transition-transform duration-300 group-hover:scale-105"
                />

                <div
                  className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent 
             opacity-70 group-hover:opacity-100 transition-opacity duration-300 rounded-full"
                />

                {/* 1. Tombol disalin dari kiri dan dikecilkan */}
                <CardAction className="absolute -bottom-2 -right-2 bg-white p-1 rounded-tl-2xl rounded-tr-lg rounded-bl-lg transition-transform duration-300">
                  <div className="flex items-center justify-center w-full h-full border-4 border-blue-400 bg-blue-800 rounded-xl transition-transform duration-300 group-hover:scale-110 p-2">
                    <ArrowRight size={20} className="text-white " />
                  </div>
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

      {/* LATEST NEWS */}

      <div className="w-full mt-10 px-5">
        {/* Header Bagian Berita */}
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-3xl font-bold">Latest News</h2>
          <div className="flex gap-2">
            {/* 4. Tambahkan event onClick pada tombol */}
            <button
              onClick={() => handleScroll(-300)}
              className="p-2 bg-gray-200 rounded-md"
            >
              <ArrowLeft size={20} />
            </button>
            <button
              onClick={() => handleScroll(300)}
              className="p-2 bg-gray-800 text-white rounded-md"
            >
              <ArrowRight size={20} />
            </button>
          </div>
        </div>

        {/* Kontainer Kartu yang Bisa di-scroll */}
        <div
          // 5. Tautkan ref dan tambahkan kelas untuk menyembunyikan scrollbar
          ref={scrollContainerRef}
          className="flex gap-6 overflow-x-auto pb-4 scrollbar-hide"
        >
          {dummyNewsData.map((news, index) => (
            <Card
              key={index}
              className="group relative h-96 min-w-80 bg-slate-800 rounded-3xl p-5 overflow-visible"
            >
              {/* Gambar Latar */}
              <img
                src={news.imageUrl}
                alt={news.title}
                className="absolute inset-0 w-full h-full object-cover rounded-3xl"
              />
              {/* Gradient Overlay */}
              <div className="absolute inset-0 bg-gradient-to-t from-black/80 to-transparent rounded-3xl" />

              {/* Konten Teks */}
              <div className="relative h-full flex flex-col justify-end text-white">
                <h3 className="text-xl font-bold">{news.title}</h3>
                <p className="mt-2 text-sm text-gray-300">{news.description}</p>
                <div className="flex items-center gap-3 mt-4">
                  <img
                    src={news.authorAvatar}
                    alt={news.authorName}
                    className="w-8 h-8 rounded-full object-cover"
                  />
                  <span className="text-sm font-medium">{news.authorName}</span>
                </div>
              </div>

              {/* Tombol Aksi dengan Efek "Notch" */}
              <CardAction className="absolute -bottom-2 -right-2 bg-white p-3 rounded-tl-2xl rounded-tr-lg rounded-bl-lg transition-transform duration-300">
                <div className="flex items-center justify-center w-full h-full border-4 border-blue-400 bg-blue-800 rounded-xl transition-transform duration-300 group-hover:scale-110 p-2">
                  <ArrowUpRight size={30} className="text-white " />
                </div>
              </CardAction>
            </Card>
          ))}
        </div>
      </div>

      {/* POPULAR SECTION */}
      <div className="w-full mt-10 px-5 py-10 text-white bg-gray-800">
        {/* Header */}
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-3xl font-bold">Popular Story</h2>
          <button className="flex items-center gap-2 bg-orange-500 text-white px-4 py-2 rounded-full text-sm">
            See more news
            <ArrowUpRight size={16} />
          </button>
        </div>

        {/* Kontainer Grid */}
        <div className="grid grid-cols-3 gap-6 mt-6">
          <Card className="col-span-3 flex flex-row  rounded-3xl p-0 h-80 bg-gray-800 border-none shadow-none text-white">
            <div className="relative group w-2/5 h-full overflow-visible">
              <img
                src={mainPopularStory.imageUrl}
                alt={mainPopularStory.title}
                className="absolute inset-0 w-full h-full object-cover rounded-3xl"
              />
              <CardAction className="group absolute -bottom-1 -right-2 bg-gray-800 p-4 rounded-tl-[40px] rounded-tr-2xl rounded-bl-2xl  transition-transform duration-300">
                <div className="flex p-4 items-center justify-center w-full h-full  border-8 border-blue-400 bg-blue-800 rounded-4xl transition-transform duration-300 group-hover:scale-110">
                  <ArrowUpRight size={28} className=" text-white" />
                </div>
              </CardAction>
            </div>
            <div className="w-3/5 p-8 flex flex-col justify-center">
              <h3 className="text-3xl font-bold">{mainPopularStory.title}</h3>
              <p className="mt-4 text-gray-600">
                {mainPopularStory.description}
              </p>
              <div className="flex items-center gap-3 mt-4">
                <img
                  src={mainPopularStory.authorAvatar}
                  alt={mainPopularStory.authorName}
                  className="w-8 h-8 rounded-full object-cover"
                />
                <span className="text-sm font-medium">
                  {mainPopularStory.authorName}
                </span>
              </div>
            </div>
          </Card>
          {/* === KARTU SEKUNDER === */}
          {secondaryPopularStories.map((story, index) => (
            <Card
              key={index}
              className="col-span-1 flex flex-row bg-gray-700 shadow-none border-none text-white rounded-3xl p-4 gap-4 h-full "
            >
              <Card className="relative group w-[50%] h-full rounded-2xl overflow-visible">
                <img
                  src={story.imageUrl}
                  alt={story.title}
                  className="absolute inset-0 w-full h-full object-cover rounded-2xl"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent rounded-2xl" />

                <CardAction className="absolute -bottom-1 -right-1 bg-gray-700 p-1 rounded-tl-2xl rounded-tr-lg rounded-bl-lg">
                  <div className="flex items-center justify-center w-full h-full border-4 border-blue-400 bg-blue-800 rounded-xl group-hover:scale-110 transition-transform duration-300 p-1">
                    <ArrowRight size={20} className="text-white" />
                  </div>
                </CardAction>
              </Card>

              <CardHeader className="flex flex-col justify-center items-start w-[50%] p-2 pl-4">
                <CardTitle className="text-base font-bold">
                  {story.title}
                </CardTitle>
                <div className="flex items-center gap-2 mt-auto pt-2">
                  <img
                    src={story.authorAvatar}
                    alt={story.authorName}
                    className="w-6 h-6 rounded-full object-cover"
                  />
                  <span className="text-xs text-gray-400">
                    {story.authorName}
                  </span>
                </div>
              </CardHeader>
            </Card>
          ))}
        </div>
      </div>

      {/* HIGHLIGHT SECTION */}
      <div className="w-full mt-10 px-5 py-10 ">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-3xl font-bold">Highlight</h2>
          <button className="flex items-center gap-2 bg-blue-500 text-white px-4 py-2 rounded-full text-sm">
            See more news
            <ArrowUpRight size={16} />
          </button>
        </div>

        {/* Kontainer Grid */}
        <div className="grid grid-cols-4 gap-6 mt-6">
          <Card className="rounded-4xl col-span-4 group relative w-full h-[550px] p-8 overflow-visible transition-all duration-300 ease-in-out">
            <img
              alt="Background"
              src={mainHighlightStory.imageUrl}
              className="absolute inset-0 w-full h-full object-cover rounded-4xl"
            />
            <div
              className="absolute inset-0 bg-gradient-to-t from-black/80 to-transparent
               opacity-100 group-hover:opacity-100 transition-opacity duration-300 rounded-4xl"
            />
            <div className="relative flex flex-col h-full justify-end text-white max-w-[75%]">
              <h2 className="text-2xl font-bold">{mainHighlightStory.title}</h2>
              <p className="mt-2 text-gray-300">
                {mainHighlightStory.description}
              </p>
              <div className="flex items-center gap-3 mt-4">
                <img
                  src={mainHighlightStory.authorAvatar}
                  alt={mainHighlightStory.authorName}
                  className="w-8 h-8 rounded-full"
                />
                <p className="text-gray-300">{mainHighlightStory.authorName}</p>
              </div>
            </div>
            <CardAction className="group absolute -bottom-2 -right-3 bg-white p-5 rounded-tl-[40px] rounded-tr-2xl rounded-bl-2xl  transition-transform duration-300">
              <div className="flex p-4 items-center justify-center w-full h-full  border-8 border-blue-400 bg-blue-800 rounded-4xl transition-transform duration-300 group-hover:scale-110">
                <ArrowRight size={64} className=" text-white" />
              </div>
            </CardAction>
          </Card>

          {/* === KARTU SEKUNDER === */}
          {dummyNewsData.map((news, index) => (
            <Card
              key={index}
              className="group relative h-90 min-w-68 bg-slate-800 rounded-3xl p-5  overflow-visible"
            >
              {/* Gambar Latar */}
              <img
                src={news.imageUrl}
                alt={news.title}
                className="absolute inset-0 w-full h-full object-cover rounded-3xl"
              />
              {/* Gradient Overlay */}
              <div className="absolute inset-0 bg-gradient-to-t from-black/80 to-transparent rounded-3xl" />

              {/* Konten Teks */}
              <div className="relative h-full flex flex-col justify-end text-white">
                <h3 className="text-xl font-bold">{news.title}</h3>
                <p className="mt-2 text-sm text-gray-300">{news.description}</p>
                <div className="flex items-center gap-3 mt-4">
                  <img
                    src={news.authorAvatar}
                    alt={news.authorName}
                    className="w-8 h-8 rounded-full object-cover"
                  />
                  <span className="text-sm font-medium">{news.authorName}</span>
                </div>
              </div>

              {/* Tombol Aksi dengan Efek "Notch" */}
              <CardAction className="absolute -bottom-2 -right-2 bg-white p-3 rounded-tl-2xl rounded-tr-lg rounded-bl-lg transition-transform duration-300">
                <div className="flex items-center justify-center w-full h-full border-4 border-blue-400 bg-blue-800 rounded-xl transition-transform duration-300 group-hover:scale-110 p-2">
                  <ArrowUpRight size={30} className="text-white " />
                </div>
              </CardAction>
            </Card>
          ))}
        </div>
      </div>

      {/* BEFORE FOOTER SECTION */}
      <div className="my-10 flex justify-center items-center w-full text-4xl font-extrabold">
        LETS CONNECT WITH DEVELAPAR!
      </div>
    </div>
  );
};

export default Homepage;
