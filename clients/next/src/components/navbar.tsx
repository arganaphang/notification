"use client";
import * as React from "react";
import { BellIcon } from "@heroicons/react/24/outline";

export default function Navbar() {
  const [count, setCount] = React.useState(1);

  return (
    <div className="fixed top-0 left-0 right-0 h-20 flex px-12 md:px-24 py-4 justify-between items-center shadow-sm">
      <h1 className="font-bold tracking-wide">Logo</h1>
      <ul>
        <div className="relative p-2 flex justify-center items-center rounded-full bg-slate-200">
          <BellIcon className="w-6 h-6 stroke-slate-400" />
          <span className="absolute bottom-0 right-0 translate-x-[6px] translate-y-1 text-[.5rem] rounded-full border-2 border-white bg-green-400 px-[6px] py-[2px] text-green-50">
            {count}
          </span>
        </div>
      </ul>
    </div>
  );
}
