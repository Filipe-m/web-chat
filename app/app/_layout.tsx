import { useFonts } from "expo-font";
import { Slot } from "expo-router";
import * as SplashScreen from "expo-splash-screen";
import { StatusBar } from "expo-status-bar";
import { useEffect } from "react";
import "react-native-reanimated";

import React from "react";
import { SessionProvider } from "@/contexts/authContext";
import { UserProvider } from "@/contexts/userContext";

SplashScreen.preventAutoHideAsync();

export default function RootLayout() {
  const [loaded] = useFonts({
    SpaceMono: require("../assets/fonts/SpaceMono-Regular.ttf"),
  });

  useEffect(() => {
    if (loaded) {
      SplashScreen.hideAsync();
    }
  }, [loaded]);

  if (!loaded) {
    return null;
  }

  return (
    <>
      <UserProvider>
        <SessionProvider>
          <Slot screenOptions={{ headerShown: false }}></Slot>
          <StatusBar style="auto" />
        </SessionProvider>
      </UserProvider>
    </>
  );
}
