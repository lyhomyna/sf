import { useContext } from "react";

import Button from "./Button.jsx";
import { AuthContext, DirItemsContext } from "storage/SfContext.jsx";
import { authServiceBaseUrl, fileServiceBaseUrl } from "config/constants.js";
import { useUploadFile } from "hooks/useUploadFile.js";
import { useLocation, useNavigate } from "react-router-dom";

export default function ControlButtons({...props}) {
    const { changeAuthStatus } = useContext(AuthContext);
    const { addDirItems } = useContext(DirItemsContext);
    const { uploadFile } = useUploadFile();
    const location = useLocation();
    const navigate = useNavigate();

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
	    const uploadPromises = fileOrFiles.map(file => uploadFile(file));
	    await Promise.all(uploadPromises);
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
		input.showPicker({ multiple: multiple });
	    } else {
		input.click();
	    }
	});

	return await fileOrFiles;
    }

    const logout = async () => {
	const res = await fetch(`${authServiceBaseUrl}/logout`, {
	    method: "GET",
	    credentials: "include", 
	})
	
	if (res.status !== 200) {
	    console.log(await res.json());
	    return;
	}
	
	changeAuthStatus();
    }

    const uploadFolder = async () => {
	let folderName = prompt("Enter new folder's name")
	if (!folderName || folderName.trim() === "") {
	    return;
	}

	try {
	    const formData = new FormData();
	    formData.append("name", folderName)
	    formData.append("curr_dir", location.pathname)

	    const res = await fetch(`${fileServiceBaseUrl}/create-directory`, {
		method: "POST",
		body: formData,
		credentials: "include",
	    })

	    if (res.status === 204) {
		alert("Wrong parent directory. The page will return to the root folder")
		navigate("/");
		return;
	    }

	    if (res.status === 409) {
		alert(`Directory '${folderName}' already exist`);
		return;
	    }
	    
	    const resJson = await res.json();
	    if (res.status === 500) {
		console.error(resJson.data);
		alert("Refresh the page and try again");
		return;
	    }

	    if (res.status === 200) {
		addDirItems({
		    dirItems: [ 
			{
			    id: resJson.id,
			    name: resJson.name,
			    path: resJson.fullPath,
			    type: "dir"
			}
		    ]
		});
	    }

	} catch(e) {
	    console.log(e);
	}
    }

    return (<div {...props}>
	<Button className="bg-neutral-700" text="Log out" onClick={logout}/>
	<Button className="bg-neutral-700" text="Upload file" onClick={ uploadFileOrFiles } />
	<Button className="bg-neutral-700" text="Create Dir" onClick={ uploadFolder }/>
    </div>);}
