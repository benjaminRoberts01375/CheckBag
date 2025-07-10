import { useContext } from "react";
import { Context, ContextType } from "./context-object";

// Custom hook to use the context
export const useList = (): ContextType => {
	const context = useContext(Context);
	if (context === undefined) {
		throw new Error("useLists must be used within a ContextProvider");
	}
	return context;
};
