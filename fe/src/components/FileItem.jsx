import Button from "./Button.jsx";

export default function FileItem({ fullFilename }) {
    const [_, ext] = fullFilename.split(".")

    const deleteItem = async () => {
	try {
	    const res = await fetch(`http://localhost:8080/delete/${fullFilename}`, {
		method: "DELETE", 
	    });
	    console.log(res)

	} catch (err) {
	    console.error("Error deleting file:", err)
	    alert("An error occured. Try again.")
	}
    }

    const downloadItem = (e) => {
	// TODO
    }

    return (<li key={filename} className="flex flex-row justify-between gap-2 items-center mt-3">
	<div className="flex gap-[0.3rem]">
	    <button onClick={deleteItem} className="w-[15px] h-[2.1rem] bg-red-200 hover:bg-red-700 duration-300" title="Delete" />
	    <div className="flex flex-row gap-x-2">
		<div className="border border-stone-300 text-slate-300 p-1">
		    .{ ext }
		</div>
		<p className="text-xl text-slate-300" title={fullFilename}>
		    { fullFilename }
		</p>
	    </div>
	</div>
	<Button className="bg-neutral-200" text="Download"/>
    </li>);
}
