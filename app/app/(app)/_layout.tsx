import { Text } from "react-native";
import { Redirect, Stack } from "expo-router";

import { useSession } from "@/contexts/authContext";
import React from "react";
import { StatusBar } from "expo-status-bar";

export default function AppLayout() {
  const { session, isLoading } = useSession();

  if (isLoading) {
    return <Text>Loading...</Text>;
  }

  if (!session) {
    return <Redirect href="/login" />;
  }

  return (
    <>
      <Stack screenOptions={{ headerShown: false }}></Stack>
      <StatusBar style="auto" />
    </>
  );
}
