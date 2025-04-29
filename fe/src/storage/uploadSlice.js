import { createSlice } from "@reduxjs/toolkit"

const initialState = {
    uploading: [],
};

const uploadSlice = createSlice({
    name: "upload",
    initialState,
    reducers: {
	addUpload: (state, action) => {
	    state.uploading.push(action.payload);
	},
	removeUpload: (state, action) => {
	    state.uploading = state.uploading.filter(f => f.id !== action.payload.id);
	},
    },
});

export const { addUpload, removeUpload } = uploadSlice.actions;
export default uploadSlice.reducer;
