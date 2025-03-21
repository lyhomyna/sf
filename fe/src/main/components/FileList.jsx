import { useEffect, useContext } from 'react';

import { FilesContext } from "../../storage/SfContext.jsx";
import { fileServiceBaseUrl } from '../../App.jsx';

import FileItem from "./FileItem.jsx";

export default function FileList({filenames}) {
    const { addFilenames } = useContext(FilesContext);

    useEffect(() => {
	fetch(`${fileServiceBaseUrl}/filenames`)
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

    return filenames.length === 0 ? (
	    <p className="text-stone-100 text-xl text-center w-[40rem]">No files uploaded yet.</p>
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
