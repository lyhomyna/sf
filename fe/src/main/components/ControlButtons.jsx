import { useContext } from "react";

import Button from "./Button.jsx";
import { AuthContext } from "storage/SfContext.jsx";
import { authServiceBaseUrl } from "config/constants.js";
import { useUploadFile } from "hooks/useUploadFile.js";

export default function ControlButtons({...props}) {
    const { changeAuthStatus } = useContext(AuthContext);
    const { uploadFile } = useUploadFile();

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
	const res = await fetch(`${authServiceBaseUrl}/logout`)
	
	if (res.status !== 200) {
	    console.log(await res.json());
	    return;
	}
	
	changeAuthStatus();
    }

    return (<div {...props}>
	<Button className="bg-neutral-700" text="Log out" onClick={logout}/>
	<Button className="bg-neutral-700" text="Upload file" onClick={ uploadFileOrFiles } />
	<Button className="bg-neutral-700" text="Create Dir" />
    </div>);}
