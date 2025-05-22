export default function DirItem({ dir }) {
    // TODO: in dir object should be path variable
    return <a href={dir.path} className="flex flex-row gap-3 mt-3 hover:bg-gray-500 rounded">
	<img src="/images/dir.png" className="h-[2.1rem]"/>
	<p className="text-xl text-slate-300">{ dir.name }</p>
    </a>;
}
