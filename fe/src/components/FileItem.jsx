import { useContext } from "react"; 
import { FilesContext } from "../storage/FilesContext.jsx";

import Button from "./Button.jsx";

export default function FileItem({ fullFilename }) {
    const { deleteFilename } = useContext(FilesContext);

    const [_, ext] = fullFilename.split(".");

    const deleteItem = async () => {
	const res = await fetch(`http://localhost:8080/delete/${fullFilename}`, {
	    method: "DELETE",
	});

	if (res.status === 204) {
	    alert("File has already deleted.");
	} else if (res.status !== 200) {
	    alert("Something went wrong while deleting file. Try again.");
	}
	
	// allert("File deleted successfuly");
	deleteFilename(fullFilename);
    }

    return (<div className="flex flex-row justify-between gap-2 items-center mt-3">
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
    </div>);
}
