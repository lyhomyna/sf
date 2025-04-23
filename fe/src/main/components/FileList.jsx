import { useEffect, useContext } from "react";

import { FilesContext } from "storage/SfContext.jsx";
import { fileServiceBaseUrl } from "config/constants.js";
import FileItem from "./FileItem.jsx";

export default function FileList() {
    const { filenames, addFilenames } = useContext(FilesContext);

    // fetch filenames on the first page load 
    useEffect(() => {
	(async () => {
	    const res = await fetch(`${fileServiceBaseUrl}/filenames`, {
		credentials: "include",
	    })

	    const resJson = await res.json()
	    if (res.status !== 200) {
		console.error(resJson.data)
		return
	    }

            addFilenames({ filenames: [...resJson.data], rewrite: true });
	})()
    }, [])

    return filenames.length === 0 ? (
	    <p className="text-stone-100 text-xl text-center w-[29rem] mt-2">No files uploaded yet.</p>
	) : (
	    <ul className="flex flex-col justify-start w-max">
	    { 
		filenames.map((filename) => {
		    return (<li key={filename}>
			<FileItem fullFilename={filename}/>
		    </li>);
		})
	    }
	    </ul>
	);
}
