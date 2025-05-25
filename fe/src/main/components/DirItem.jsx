import { useContext } from "react"; 
import { DirItemsContext } from "storage/SfContext.jsx";
import { fileServiceBaseUrl } from "config/constants.js";

export default function DirItem({ dir }) {
    const { deleteDirItem } = useContext(DirItemsContext);

    const deleteDir = async () => {
	try {
	    const res = await fetch(`${fileServiceBaseUrl}/delete-directory/${dir.id}`, {
		method: "DELETE",
	    });

	    if (!res.ok) {
		const resJson = await res.json();
		console.error(resJson.data);
	    } 	

	    deleteDirItem(dir);
	} catch (err) {
	    console.error(err);
	}
    };

    return <>
	<div className="flex gap-[0.3rem]">
	    <button onClick={ deleteDir } className="w-[15px] h-[2.1rem] bg-red-200 hover:bg-red-700 duration-300" title="Delete" />
	    <a href={dir.path} className="flex flex-row gap-3 mt-3 hover:bg-gray-500 rounded">
		<img src="/images/dir.png" className="h-[2.1rem]"/>
		<p className="text-xl text-slate-300">{ dir.name }</p>
	    </a>;
	</div>
    </>;
}
