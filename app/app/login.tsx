import { useSession } from "@/contexts/authContext";
import { Link, useRouter } from "expo-router";
import React, { useEffect } from "react";
import { StyleSheet, View, Text, TextInput, Pressable } from "react-native";
import Feather from "@expo/vector-icons/Feather";

export default function Login() {
  const [email, onChangeEmail] = React.useState("");
  const [password, onChangePassword] = React.useState("");
  const [showPassword, onChangeShowPassword] = React.useState(true);
  const { session, signIn } = useSession();
  const router = useRouter();

  useEffect(() => {
    if (session) {
      router.replace("/");
    }
  }, [session]);

  const handleLogin = () => {
    try {
      signIn(email, password);
      router.replace("/");
    } catch (error) {
      console.error("Erro ao fazer login:", error);
      alert("Falha ao fazer login. Verifique suas credenciais.");
    }
  };

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
        <Pressable
          onPress={handleLogin}
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
        Não possui uma conta?{" "}
        <Link style={styles.registerLink} href="/register">
          Registre-se
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
    width: "100%",
    marginBottom: 12,
    borderWidth: 1,
    borderRadius: 10,
    padding: 10,
  },
  confirmBtn: {
    borderRadius: 8,
    padding: 10,
    alignItems: "center",
    justifyContent: "center",
    width: "100%",
    marginTop: 12,
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
  passwordInput: {
    height: 40,
    width: "100%",
    borderRadius: 10,
    paddingLeft: 10,
    paddingRight: 40,
    textAlignVertical: "center",
  },
});
