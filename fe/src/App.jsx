import TopBar from "./components/TopBar.jsx";
import FileList from "./components/FileList.jsx";

export default function App() {
    const email="supercoolemail@super.mail" 
    const imageURL="https://static.vecteezy.com/system/resources/previews/018/871/797/non_2x/happy-cat-transparent-background-png.png"

    return (<div className="p-2" >
	<TopBar email={email} imageURL={imageURL} />
	<FileList />
    </div>);
}
