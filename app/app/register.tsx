import { Link } from "expo-router";
import React from "react";
import { StyleSheet, View, Text, TextInput, Pressable } from "react-native";
import Feather from "@expo/vector-icons/Feather";

export default function Register() {
  const [email, onChangeEmail] = React.useState("");
  const [password, onChangePassword] = React.useState("");
  const [confirmPassword, onChangeConfirmPassowrd] = React.useState("");
  const [showPassword, onChangeShowPassword] = React.useState(true);
  const [showConfirmPassword, onChangeShowConfirmPassword] =
    React.useState(true);

  return (
    <View style={styles.container}>
      <View style={styles.formContainer}>
        <Text style={styles.texto}>Email:</Text>
        <TextInput
          style={styles.input}
          onChangeText={onChangeEmail}
          value={email}
          placeholder="Email"
          keyboardType="email-address"
          autoCapitalize="none"
        />

        <Text style={styles.texto}>Senha:</Text>
        <View style={styles.passwordContainer}>
          <TextInput
            style={styles.passwordInput}
            onChangeText={onChangePassword}
            value={password}
            placeholder="Senha"
            secureTextEntry={showPassword}
            autoCapitalize="none"
          />
          <Pressable
            onPress={() => {
              onChangeShowPassword(!showPassword);
            }}
            style={styles.eyeIcon}
          >
            {showPassword ? (
              <Feather name="eye" size={18} color="black" />
            ) : (
              <Feather name="eye-off" size={18} color="black" />
            )}
          </Pressable>
        </View>

        <Text style={styles.texto}>Confirme a senha:</Text>
        <View style={styles.passwordContainer}>
          <TextInput
            style={styles.passwordInput}
            onChangeText={onChangeConfirmPassowrd}
            value={confirmPassword}
            placeholder="Senha"
            secureTextEntry={showConfirmPassword}
            autoCapitalize="none"
          />
          <Pressable
            onPress={() => {
              onChangeShowConfirmPassword(!showConfirmPassword);
            }}
            style={styles.eyeIcon}
          >
            {showConfirmPassword ? (
              <Feather name="eye" size={18} color="black" />
            ) : (
              <Feather name="eye-off" size={18} color="black" />
            )}
          </Pressable>
        </View>

        <Pressable
          onPress={() => {}}
          style={({ pressed }) => [
            {
              backgroundColor: pressed ? "#b39cd0" : "#845ec2",
            },
            styles.confirmBtn,
          ]}
        >
          {({ pressed }) => (
            <Text style={{ color: pressed ? "#000000" : "#FFFFFF" }}>
              Confirmar
            </Text>
          )}
        </Pressable>
      </View>

      <Text style={styles.registerText}>
        Já possui uma conta?{" "}
        <Link style={styles.registerLink} href="/">
          Login
        </Link>
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  },
  formContainer: {
    width: "80%",
    alignItems: "center",
  },
  texto: {
    marginBottom: 4,
    alignSelf: "flex-start",
  },
  input: {
    height: 40,
    borderWidth: 1,
    width: "100%",
    marginBottom: 12,
    borderRadius: 10,
    paddingLeft: 10,
    paddingRight: 40,
    textAlignVertical: "center",
  },
  passwordInput: {
    height: 40,
    width: "100%",
    borderRadius: 10,
    paddingLeft: 10,
    paddingRight: 40,
    textAlignVertical: "center",
  },
  confirmBtn: {
    borderRadius: 8,
    padding: 10,
    alignItems: "center",
    justifyContent: "center",
    width: "90%",
    marginTop: 12,
  },
  passwordContainer: {
    height: 40,
    width: "100%",
    marginBottom: 12,
    borderWidth: 1,
    borderRadius: 10,
    flexDirection: "row",
    alignItems: "center",
    position: "relative",
  },
  eyeIcon: {
    position: "absolute",
    right: 10,
    top: "50%",
    transform: [{ translateY: -9 }],
  },
  registerText: {
    position: "absolute",
    bottom: 40,
    fontSize: 14,
    textAlign: "center",
  },
  registerLink: {
    color: "#007BFF",
  },
});
