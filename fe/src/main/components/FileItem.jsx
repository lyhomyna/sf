import { useContext } from "react"; 
import { FilesContext } from "storage/SfContext.jsx";
import { fileServiceBaseUrl } from "config/constants.js";

import Button from "./Button.jsx";

export default function FileItem({ file }) {
    const { deleteFile } = useContext(FilesContext);

    const deleteItem = async () => {
	const res = await fetch(`${fileServiceBaseUrl}/delete/${file.id}`, {
	    method: "DELETE",
	});

	if (res.status === 204) {
	    alert("File has already deleted.");
	} else if (res.status !== 200) {
	    alert("Something went wrong while deleting file. Try again.");
	    return
	}
	
	// allert("File deleted successfuly");
	deleteFile(file);
    }

    const downloadFile = async () => {
	try {
	    const res = await fetch(`${fileServiceBaseUrl}/download/${file.id}`);
	    if (res.status === 404) {
		alert("File not found. Refresh the page.");
		return;
	    } else if (!res.ok) {
		console.error("Download failed:", res.statusText);
                return;
            }

	    const blob = await res.blob();
	    const url = window.URL.createObjectURL(blob);

	    const a = document.createElement("a");
	    a.href = url;
	    a.download = file.filename;

	    document.body.appendChild(a);
	    a.click();
	    document.body.removeChild(a);

	    window.URL.revokeObjectURL(url);
	} catch (e) {
	    console.error(e)
	    alert("Something went wrong downloading the file. Refresh the page and try again.")
	}
    }
    
    return <li>
	<div className="flex flex-row justify-between gap-2 items-center mt-3">
	    <div className="flex gap-[0.3rem]">
		<button onClick={deleteItem} className="w-[15px] h-[2.1rem] bg-red-200 hover:bg-red-700 duration-300" title="Delete" />
		<div className="flex flex-row gap-x-2">
		    <div className="border border-stone-300 text-slate-300 p-1">
			.{ file.filename.split(".")[1] }
		    </div>
		    <p className="text-xl text-slate-300" title={file.filename}>
			{ file.filename }
		    </p>
		</div>
	    </div>
	    <Button className="bg-neutral-700" text="Download" onClick={downloadFile}/>
	</div>
    </li>;
}
