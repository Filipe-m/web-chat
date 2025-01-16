import { PropsWithChildren } from "react";
import { Pressable, Text, StyleSheet, ViewStyle, TextStyle } from "react-native";

type DefaultBtnProps = PropsWithChildren<{
  method: () => void;
  text: string;
  containerStyle?: ViewStyle; // Estilo do botão no estado normal
  textStyle?: TextStyle;      // Estilo do texto no estado normal
  pressedContainerStyle?: ViewStyle; // Estilo do botão no estado pressionado
  pressedTextStyle?: TextStyle;      // Estilo do texto no estado pressionado
}>;

export default function DefaultBtn({
  text,
  method,
  containerStyle,
  textStyle,
  pressedContainerStyle,
  pressedTextStyle,
}: DefaultBtnProps) {
  return (
    <Pressable
      onPress={method}
      style={({ pressed }) => [
        styles.defaultContainer,
        pressed ? pressedContainerStyle : containerStyle,
      ]}
    >
      {({ pressed }) => (
        <Text
          style={[
            styles.defaultText,
            pressed ? pressedTextStyle : textStyle,
          ]}
        >
          {text}
        </Text>
      )}
    </Pressable>
  );
}

const styles = StyleSheet.create({
  defaultContainer: {
    borderRadius: 8,
    padding: 10,
    alignItems: "center",
    justifyContent: "center",
    width: "100%",
  },
  defaultText: {
    fontSize: 16,
    fontWeight: "bold",
    color: "#FFFFFF",
  },
});
