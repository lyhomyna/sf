import ControlButtons from "./ControlButtons.jsx";

export default function TopBar({ email, imageURL }) {
    return (<div className="flex flex-row items-center gap-x-5">
	<div className="w-16 h-16 bg-stone-700 rounded-md shadow-md">
	    <img src={imageURL} alt="Avatar"/>
	</div>
	<div className="font-medium text-slate-300">
	    { email }
	</div>
	<ControlButtons className="flex gap-x-1" />
    </div>);
}
