import { useState } from "react";
import { Provider } from "react-redux";
import { Routes, Route } from 'react-router-dom';

import TopBar from "./components/TopBar.jsx";
import FileList from "./components/FileList.jsx";
import DragAndDrop from "./components/DragAndDrop.jsx";
import { DirItemsContext } from "storage/SfContext.jsx";
import store from "storage/store.js"; 

export default function MainPage() {
    const [dirItems, setDirItems] = useState([]);

    const addDirItems = ({ dirItems, rewrite=false }) => {
	// rewrite option is used only on first files load 
	if (rewrite) {
	    setDirItems(dirItems);
	} else {
	    setDirItems(currDirItems => [...currDirItems, ...dirItems]);
	}
    };

    const deleteDirItem = (file) => {
	setDirItems(oldFiles => { 
	    return oldFiles.filter(f => f.id !== file.id)
	});
    };

    return <Provider store={store}>
	    <DirItemsContext value={ { 
		dirItems: dirItems, 
		addDirItems: addDirItems, 
		deleteDirItem: deleteDirItem, 
	    }}>
		<div className="p-2" >
		    <TopBar />
		    <Routes>
			<Route path="*" element={<FileList />} />
		    </Routes>
		    <DragAndDrop /> 
		</div>
	    </DirItemsContext>
	</Provider>
}
