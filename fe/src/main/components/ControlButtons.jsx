import { useContext } from "react";

import Button from "./Button.jsx";
import { AuthContext, FilesContext } from "../../storage/SfContext.jsx";
import { authServiceBaseUrl, fileServiceBaseUrl } from "../../App.jsx";

export default function ControlButtons({...props}) {
    const { addFilename } = useContext(FilesContext);
    const { changeAuthStatus } = useContext(AuthContext);

    const uploadFileOrFiles = async (multiple = false) => {
	const supportsFileSystemAccess = "showOpenFilePicker" in window &&
	    (() => {
		try {
		    return window.self === window.top;
		} catch {
		    return false;
		}
	    })();
	
	let fileOrFiles = undefined;

	if (supportsFileSystemAccess) {
	    fileOrFiles = await openFilePicker(multiple);
	} else {
	    fileOrFiles = await openStandardFilePicker(multiple);
	}
	
	if (fileOrFiles) {
	    await fileOrFiles.forEach(async (file) => {
		const formData = new FormData();
		formData.append("file", file);
		try {
		    const response = await fetch(`${fileServiceBaseUrl}/save`, {
			method: "POST", 
			body: formData,
			credentials: "include",
		    });
		    
		    if (response.ok) {
			// alert("File uploaded successfully!");
			addFilename(file.name) // to show filenames in list
		    } else if(response.status === 400) {
			alert("File with the same filename has already uploaded.")
		    } else {
			alert("Failed to upload file. Try again.");
		    }
		} catch (err) {
		    console.error("Error uploading file:", err)
		    alert("An error occured. Try again.")
		}
	    });
	}
    }

    const openFilePicker = async (multiple) => {
	let fileOrFiles = undefined;	
	try {
	    const handles = await window.showOpenFilePicker({ multiple });
	    if (!multiple) {
		fileOrFiles = await handles[0].getFile();
		fileOrFiles.handle = handles[0];
	    } else {
		fileOrFiles = await Promise.all(
		    handles.map(async (handle) => {
			const file = await handle.getFile();
			file.handle = handle;
			return file;
		    })
		);
	    }
	} catch (err) {
	    if (err.name !== 'AbortError') {
		console.error(err.name, err.message);
	    }
	}
	return fileOrFiles;
    }

    const openStandardFilePicker = async (multiple) => {
	const fileOrFiles = new Promise((resolve) => {
	    const input = document.createElement("input");
	    input.type = "file";
	    input.style.display = "none";

	    if (multiple) {
		input.multiple = true;
	    }

	    document.body.append(input);

	    input.addEventListener("change", () => {
		// remove input from DOM
		input.remove();

		if (!input.files) {
		    return;
		}
		
		// return all files or just one from promise
		resolve(multiple ? Array.from(input.files) : input.files[0]);
	    });

	    if ('showPicker' in HTMLInputElement.prototype) {
		input.showPicker();
	    } else {
		input.click();
	    }
	});

	return await fileOrFiles;
    }

    const logout = async () => {
	const res = await fetch(`${authServiceBaseUrl}/logout`)
	
	if (res.status !== 200) {
	    console.log(await res.json());
	    return;
	}
	
	changeAuthStatus();
    }

    return (<div {...props}>
	<Button className="bg-neutral-700" text="Log out" onClick={logout}/>
	<Button className="bg-neutral-300" text="Cnange password" />
	<Button className="bg-neutral-700"text="Upload file" onClick={uploadFileOrFiles} />
    </div>);
}
