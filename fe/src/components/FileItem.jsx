export default function FileItem({ fullFilename }) {
    const [filename, ext] = fullFilename.split(".")
    return (<div className="flex flex-row gap-x-3">
	<div className="border border-black">
	    .{ ext }
	</div>
	<p>
	    {filename}
	</p>
    </div>);
}
