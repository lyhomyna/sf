import { useContext } from "react";
import { FilesContext } from "storage/SfContext.jsx";
import { fileServiceBaseUrl } from "config/constants.js";

export function useUploadFile() {
    const { addFilename } = useContext(FilesContext);

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
		addFilename(file.name) // to show filenames in list
	    } else if(response.status === 400) {
		alert(`File ${file.name} has already uploaded.`)
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
