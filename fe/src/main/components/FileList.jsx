import { useEffect, useContext } from "react";
import { useLocation } from 'react-router-dom';

import { DirItemsContext } from "storage/SfContext.jsx";
import { fileServiceBaseUrl } from "config/constants.js";
import FileItem from "./FileItem.jsx";
import DirItem from "./DirItem.jsx";

export default function FileList() {
    const location = useLocation();
    const { dirItems, addDirItems } = useContext(DirItemsContext);

   // fetch filenames on the first page load 
    useEffect(() => {
        (async () => {
	    try {
		// that's ok, because pathname always starts with '/'
		const res = await fetch(`${fileServiceBaseUrl}${location.pathname}`, {
		    credentials: "include",
		})

		const resJson = await res.json()
		if (res.status !== 200) {
		    console.error(resJson.data)
		    return
		}

		if (resJson === null) {
		    return
		}

		addDirItems({ dirItems: [...resJson], rewrite: true });
	    } catch (e) {
		console.error(e);
	    }
        })()
    }, [])

    return dirItems.length === 0 ? (
	    <p className="text-stone-100 text-xl text-center w-[29rem] mt-2">No files uploaded yet.</p>
	) : (
	    <ul className="flex flex-col justify-start w-max">
	    { 
		dirItems.map((dirItem) => {
		    let item;

		    if (dirItem.type === "dir") {
			item = <DirItem key={dirItem.id} dir={dirItem}/>;
		    } else {
			item = <FileItem key={dirItem.id} file={dirItem}/>; 
		    }

		    return item;
		})
	    }
	    </ul>
	);
}
