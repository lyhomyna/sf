import { useEffect, useState } from "react";

import TopBar from "./components/TopBar.jsx";
import FileList from "./components/FileList.jsx";
import DragAndDrop from "./components/DragAndDrop.jsx";
import { FilesContext } from "storage/SfContext.jsx";
import { authServiceBaseUrl } from "config/constants.js";

export default function MainPage() {
    const imageURL="https://static.vecteezy.com/system/resources/previews/018/871/797/non_2x/happy-cat-transparent-background-png.png"
    const [email, setEmail] = useState("supercoolemail@super.mail");
    const [filenames, setFilenames] = useState([]);

    useEffect(() => {
	const fetchUser = async () => {
	    try {
		const response = await fetch(`${authServiceBaseUrl}/get-user`);

		if (response.status === 200) {
		    const user = await response.json(); 
		    setEmail(user.email);
		} else if (response.status === 404) {
		    const err = await response.json();
		    console.log(err.message)
		}
	    } catch (err) {
		console.log("Failed to fetch user", err)
	    }
	}

	fetchUser();
    }, [])

    const addFilenames = (filenames) => {
	if (Array.isArray(filenames)) {
	    setFilenames(filenames);
	}
    };

    const addFilename = (filename) => {
	setFilenames(oldFilenames => [...oldFilenames, filename]);
    };

    const deleteFilename = (filename) => {
	setFilenames(oldFilenames => { 
	    let withoutFilename;
	    oldFilenames.forEach((_, i) => {
		if (oldFilenames[i] === filename) {
		    withoutFilename = [...oldFilenames.slice(0, i), ...oldFilenames.slice(i+1, oldFilenames.length)]
		    return
		}
	    });
	    return withoutFilename;
	});
    };
    return <FilesContext value={ { 
	filenames: filenames, 
	addFilenames: addFilenames, 
	addFilename: addFilename, 
	deleteFilename: deleteFilename, 
    }}>
	<div className="p-2" >
	    <TopBar email={email} imageURL={imageURL} />
	    <FileList filenames={filenames} />
	    <DragAndDrop /> 
	</div>
    </FilesContext>
}
