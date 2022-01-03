import { configureStore } from '@reduxjs/toolkit'
import appReducer, { IApplicationState } from '../features/applicationSlice';

export interface State {
  app: IApplicationState;
}

const store = configureStore({
  reducer: {
    app: appReducer,
  },
  preloadedState: {
  },
  middleware: (getDefaultMiddleware) => 
    getDefaultMiddleware({
      serializableCheck: false,
    })
});

export default store;
