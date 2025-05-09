import { useState } from "react";
import { Provider } from "react-redux";

import TopBar from "./components/TopBar.jsx";
import FileList from "./components/FileList.jsx";
import DragAndDrop from "./components/DragAndDrop.jsx";
import { FilesContext } from "storage/SfContext.jsx";
import store from "storage/store.js"; 

export default function MainPage() {
    const [files, setFiles] = useState([]);

    const addFiles = ({ files, rewrite=false }) => {
	// rewrite option is used only on first files load 
	if (rewrite) {
	    setFiles(files);
	} else {
	    setFiles(oldFiles => [...oldFiles, ...files]);
	}
    };

    const deleteFile = (file) => {
	setFiles(oldFiles => { 
	    return oldFiles.filter(f => f.id !== file.id)
	});
    };

    return <Provider store={store}>
	<FilesContext value={ { 
	    files: files, 
	    addFiles: addFiles, 
	    deleteFile: deleteFile, 
	}}>
	    <div className="p-2" >
		<TopBar />
		<FileList />
		<DragAndDrop /> 
	    </div>
	</FilesContext>
    </Provider>
}
