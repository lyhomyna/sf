import { useEffect, useContext } from 'react';

import { FilesContext } from "../storage/FilesContext.jsx";
import FileItem from "./FileItem.jsx";

export default function FileList() {
    const { filenames, addFilenames } = useContext(FilesContext);

    useEffect(() => {
	fetch("http://localhost:8080/filenames")
	.then((res) => {
	    return res.json();
	}).then((json) => {
	    if (Array.isArray(json.data)) {
                addFilenames(json.data);
            } else {
                console.error("Invalid data format:", json);
                addFilenames([]);
            }
	}).catch((err) => console.error("Failed to fetch filenames:", err));
    }, [])

    return <ul className="flex flex-col justify-start w-max">
	{ 
	    filenames.map((filename) => {
		return <FileItem fullFilename={filename}/>
	    })
	}
    </ul>
}
