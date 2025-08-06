import { Button } from "@/components/ui/button";
import {
  Card,
  CardAction,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

const dummyData = [
  {
    title: "Card Title 1",
    description: "This is a description for card 1",
    author: "Author 1",
  },
  {
    title: "Card Title 2",
    description: "This is a description for card 1",
    author: "Author 1",
  },
  {
    title: "Card Title 2",
    description: "This is a description for card 1",
    author: "Author 1",
  },
];

const Homepage = () => {
  return (
    <div className="mx-auto px-5">
      <div className="my-5">OH Hi! Welcome to develapar!</div>
      <div
        id="top-menu"
        className="flex flex-row gap-5 justify-between w-full h-[40vw] items-center "
      >
        <Card className=" w-full h-full max-w-[50%] bg-blue-300 rounded-4xl">
          <div className="flex-1" />

          <CardHeader className="flex flex-row items-start justify-between">
            <div className="flex flex-col gap-1">
              <CardTitle>Ini Judul COntentnya</CardTitle>
              <CardDescription>Ini Decsripsi conetenya apa</CardDescription>
              <CardDescription>Authornya disini bro</CardDescription>
            </div>

            <CardAction className="bg-white">
              <Button variant="ghost">Card Action</Button>
            </CardAction>
          </CardHeader>
        </Card>

        {/* Kanan */}
        <div className="w-full h-full max-w-[50%] flex flex-col gap-4 justify-between">
          {dummyData.map((item, index) => (
            <Card
              key={index}
              className=" py-2 shadow-none border-none  flex flex-row h-1/3"
            >
              <Card className=" bg-amber-700 w-[50%] h-[100%] rounded-4xl">
                <p>Ini Isinya foto ntar</p>
              </Card>

              <CardHeader className="flex px-0 flex-col justify-evenly items-start  w-full">
                <CardTitle className="pt-1">{item.title}</CardTitle>
                <CardDescription>Ini Deskripsi Card Lainnya</CardDescription>
                <CardDescription className="pb-1">
                  Ini Deskripsi Card Lainnya
                </CardDescription>
              </CardHeader>
            </Card>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Homepage;
