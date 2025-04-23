import { createSlice } from "@reduxjs/toolkit"

const initialState = {
    uploading: [],
};

const uploadSlice = createSlice({
    name: "upload",
    initialState,
    reducers: {
	addUpload: (state, action) => {
	    console.log("+1 file to upload")
	    state.uploading.push(action.payload);

	    console.log("total files to upload:", state.uploading.length)
	},
	removeUpload: (state, action) => {
	    console.log("-1 file to upload")
	    state.uploading = state.uploading.filter(f => f.id !== action.payload.id);

	    console.log("total files to upload:", state.uploading.length)
	},
    },
});

export const { addUpload, removeUpload } = uploadSlice.actions;
export default uploadSlice.reducer;
