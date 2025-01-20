export default function FileItem({ fullFilename }) {
    const [filename, ext] = fullFilename.split(".")
    return (<li className="flex flex-row gap-[0.3rem] items-center mt-3">
	<button className="w-[15px] h-[2.1rem] bg-red-400" title="Delete" />
	<div className="flex flex-row gap-x-2">
	    <div className="border border-stone-300 text-slate-300 p-1">
		.{ ext }
	    </div>
	    <p className="text-xl text-slate-300" title={fullFilename}>
		{ fullFilename }
	    </p>
	</div>
    </li>);
}
