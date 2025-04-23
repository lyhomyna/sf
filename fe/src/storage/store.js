import { configureStore } from "@reduxjs/toolkit";
import uploadReducer from "./uploadSlice.js";

const store = configureStore({
    reducer: {
	upload: uploadReducer,
    }
});

export default store;
