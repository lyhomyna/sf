import { useEffect, useContext } from "react";

import { FilesContext } from "storage/SfContext.jsx";
import { fileServiceBaseUrl } from "config/constants.js";
import FileItem from "./FileItem.jsx";

export default function FileList() {
    const { files, addFiles } = useContext(FilesContext);

    // fetch filenames on the first page load 
    useEffect(() => {
	(async () => {
	    const res = await fetch(`${fileServiceBaseUrl}/files`, {
		credentials: "include",
	    })

	    const resJson = await res.json()
	    if (res.status !== 200) {
		console.error(resJson.data)
		return
	    }

            addFiles({ files: [...resJson.data], rewrite: true });
	})()
    }, [])

    return files.length === 0 ? (
	    <p className="text-stone-100 text-xl text-center w-[29rem] mt-2">No files uploaded yet.</p>
	) : (
	    <ul className="flex flex-col justify-start w-max">
	    { 
		files.map((file) => {
		    return ( <FileItem key={file.id} file={file}/> );
		})
	    }
	    </ul>
	);
}
