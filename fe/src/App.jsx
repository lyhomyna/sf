import TopBar from "./components/TopBar.jsx";
import FileItem from "./components/FileItem.jsx";

export default function App() {
    const email="supercoolemail@super.mail" 
    const imageURL="https://static.vecteezy.com/system/resources/previews/018/871/797/non_2x/happy-cat-transparent-background-png.png"

    return (<div className="p-2" >
	<TopBar email={email} imageURL={imageURL} />
	<ul className="flex flex-col justify-start w-max">
	    <FileItem fullFilename="supercoolfile1.txt"/>
	    <FileItem fullFilename="supercoolfilee2.txt"/>
	    <FileItem fullFilename="supercoolfileee3.txt"/>
	</ul>
    </div>);
}
