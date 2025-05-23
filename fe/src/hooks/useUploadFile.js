import { useContext } from "react";
import { DirItemsContext } from "storage/SfContext.jsx";
import { fileServiceBaseUrl } from "config/constants.js";
import { useDispatch } from "react-redux";
import { addUpload, removeUpload } from "storage/uploadSlice.js";
import { useLocation } from "react-router-dom" 

export function useUploadFile() {
    const { addDirItems } = useContext(DirItemsContext);
    const dispatch = useDispatch();
    const location = useLocation();

    const uploadFile = async (file) => {
	const formData = new FormData();
	formData.append("file", file);
	formData.append("dir", location.pathname);

	// to show progress at the top right corner
	const tempFileId = file.name+randomInt();
	dispatch(addUpload({id: tempFileId}));

	try {
	    const response = await fetch(`${fileServiceBaseUrl}/save`, {
		method: "POST", 
		body: formData,
		credentials: "include",
	    });
	    
	    const resJson = await response.json();
	    if (response.status === 200) {
		// show file in list
		addDirItems({
		    dirItems: [
			{
			    id: resJson.id,
			    name: resJson.filename,
			    type: "file",
			}
		    ]
		})
	    } else if(response.status === 400) {
		alert(resJson.data)
		console.error(resJson.data)
	    } else {
		alert(`Failed to upload file '${file.name}'. Try again.`);
		console.error(resJson.data)
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
