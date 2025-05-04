import { useContext } from "react";
import { FilesContext } from "storage/SfContext.jsx";
import { fileServiceBaseUrl } from "config/constants.js";
import { useDispatch } from "react-redux";
import { addUpload, removeUpload } from "storage/uploadSlice.js";

export function useUploadFile() {
    const { addFiles } = useContext(FilesContext);
    const dispatch = useDispatch();

    const uploadFile = async (file) => {
	const formData = new FormData();
	formData.append("file", file);

	// to show progress at top right corner
	const tempFileId = file.name+randomInt();
	dispatch(addUpload({id: tempFileId}));

	try {
	    const response = await fetch(`${fileServiceBaseUrl}/save`, {
		method: "POST", 
		body: formData,
		credentials: "include",
	    });
	    
	    const resJson = await response.json();
	    if (response.ok) {
		// show file in list
		addFiles({files: [{
		    id: resJson.data.id,
		    filename: resJson.data.filename,
		    createdAt: -1, // no creation time
		}]})
	    } else if(response.status === 400) {
		alert(await resJson.data)
	    } else {
		alert(`Failed to upload file '${file.name}'. Try again.`);
		console.log(resJson.data)
	    }
	} catch (err) {
	    console.error("Error uploading file:", err)
	    alert(`Failed to upload file '${file.name}. Try again.`)
	} finally {
	    // to remove progress at top right corner
	    dispatch(removeUpload({id: tempFileId}))
	}
    };

    return { uploadFile };
}

function randomInt(min=0, max=9999999999) {
    return Math.floor(Math.random() * (max - min + 1) + min);
}
