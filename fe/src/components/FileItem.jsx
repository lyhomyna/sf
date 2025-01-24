import Button from "./Button.jsx";

export default function FileItem({ fullFilename }) {
    const [filename, ext] = fullFilename.split(".")

    return (<li key={filename} className="flex flex-row justify-between gap-2 items-center mt-3">
	<div className="flex gap-[0.3rem]">
	    <button className="w-[15px] h-[2.1rem] bg-red-400 hover:bg-red-700 duration-300" title="Delete" />
	    <div className="flex flex-row gap-x-2">
		<div className="border border-stone-300 text-slate-300 p-1">
		    .{ ext }
		</div>
		<p className="text-xl text-slate-300" title={fullFilename}>
		    { fullFilename }
		</p>
	    </div>
	</div>
	<Button text="Download"/>
    </li>);
}
