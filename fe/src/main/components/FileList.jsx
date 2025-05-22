import { useEffect, useContext } from "react";
import { useLocation } from 'react-router-dom';

import { FilesContext } from "storage/SfContext.jsx";
import { fileServiceBaseUrl } from "config/constants.js";
import FileItem from "./FileItem.jsx";
import DirItem from "./DirItem.jsx";

export default function FileList() {
    const location = useLocation();
    const { files, addFiles } = useContext(FilesContext);

   // fetch filenames on the first page load 
    useEffect(() => {
        (async () => {
	    try {
		// that's ok, cause pathname always starts with /
		const res = await fetch(`${fileServiceBaseUrl}${location.pathname}`, {
		    credentials: "include",
		})

		const resJson = await res.json()
		if (res.status !== 200) {
		    console.error(resJson.data)
		    return
		}

		addFiles({ files: [...resJson.data], rewrite: true });
	    } catch (e) {
		console.error(e);
	    }
        })()
    }, [])

    return files.length === 0 ? (
	    <p className="text-stone-100 text-xl text-center w-[29rem] mt-2">No files uploaded yet.</p>
	) : (
	    <ul className="flex flex-col justify-start w-max">
	    { 
		// TODO: dirItem should be:
		// {
		//	id: string,
		//	type: string,
		//	item: { ... },
		// }
		files.map((dirItem) => {
		    let item;

		    if (dirItem.type === "dir") {
			item = <DirItem key={dirItem.id} dir={dirItem.item}/>;
		    } else {
			item = <FileItem key={dirItem.id} file={dirItem.item}/>; 
		    }

		    return item;
		})
	    }
	    </ul>
	);
}
