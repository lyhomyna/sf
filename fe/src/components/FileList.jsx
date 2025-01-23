import { useState, useEffect } from 'react';

import FileItem from "./FileItem.jsx";

export default function FileList() {
    const [filenames, setFilenames] = useState([]); 

    useEffect(() => {
	fetch("http://localhost:8080/filenames")
	.then((res) => {
	    return res.json();
	}).then((json) => {
	    if (Array.isArray(json.data)) {
                setFilenames(json.data);
            } else {
                console.error("Invalid data format:", json);
                setFilenames([]);
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
