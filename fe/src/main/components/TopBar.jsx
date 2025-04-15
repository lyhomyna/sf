import { CircleLoader } from "react-spinners";
import ControlButtons from "./ControlButtons.jsx";

export default function TopBar({ email, imageURL }) {
    return <div className="flex flex-row justify-between items-center">
	<div className="flex flex-row items-center gap-5 flex-wrap">
	    <div className="w-16 h-16 bg-stone-700 rounded-md shadow-md">
		<img src={imageURL} alt="Avatar"/>
	    </div>
	    <div className="font-medium text-slate-300">
		{ email }
	    </div>
	    <ControlButtons className="flex items-center gap-x-1" />
	</div>
	<div class="relative group inline-block cursor-pointer pr-5">
	    <CircleLoader size="30" color="#fff"/>
	    <span class="absolute top-full right-0 px-2 py-1 bg-black text-white text-sm rounded opacity-0 group-hover:opacity-100 pointer-events-none max-w-[200px] break-words z-50">
		Your files is uploading...
	    </span>
	</div>
    </div>;
}
