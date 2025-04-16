import { useContext } from "react";
import { FilesContext } from "storage/SfContext.jsx";
import { fileServiceBaseUrl } from "config/constants.js";

export function useUploadFile() {
    const { addFilenames } = useContext(FilesContext);

    const uploadFile = async (file) => {
	const formData = new FormData();
	formData.append("file", file);
	try {
	    const response = await fetch(`${fileServiceBaseUrl}/save`, {
		method: "POST", 
		body: formData,
		credentials: "include",
	    });
	    
	    if (response.ok) {
		console.log(`File '${file.name}' uploaded.`)
		addFilenames({filenames: [file.name]}) // to show filenames in list
	    } else if(response.status === 400) {
		const res = await response.json()
		alert(await res.data)
	    } else {
		alert(`Failed to upload file '${file.name}'. Try again.`);
	    }
	} catch (err) {
	    console.error("Error uploading file:", err)
	    alert(`Failed to upload file '${file.name}. Try again.`)
	}
    };

    return { uploadFile };
}
