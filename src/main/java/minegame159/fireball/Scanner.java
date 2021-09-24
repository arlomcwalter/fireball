package minegame159.fireball;

import java.io.IOException;
import java.io.Reader;

public class Scanner {
    private final Reader reader;
    private final StringBuilder sb = new StringBuilder();

    private char current, next;
    private int line = 1;

    public Scanner(Reader reader) {
        this.reader = reader;

        advance(false);
        advance(false);
    }

    public Token next() {
        sb.setLength(0);

        skipWhitespace();
        if (isAtEnd()) return token(TokenType.Eof);

        char c = advance();
        if (isTwice(c, '&')) return token(TokenType.And);
        if (isTwice(c, '|')) return token(TokenType.Or);
        if (isDigit(c)) return number();
        if (isAlpha(c)) return identifier();

        return switch (c) {
            case '(' -> token(TokenType.LeftParen);
            case ')' -> token(TokenType.RightParen);
            case '{' -> token(TokenType.LeftBrace);
            case '}' -> token(TokenType.RightBrace);
            case ';' -> token(TokenType.Semicolon);
            case ',' -> token(TokenType.Comma);
            case '.' -> token(TokenType.Dot);

            case '!' -> token(match('=') ? TokenType.BangEqual : TokenType.Bang);
            case '=' -> token(match('=') ? TokenType.EqualEqual : TokenType.Equal);
            case '>' -> token(match('=') ? TokenType.GreaterEqual : TokenType.Greater);
            case '<' -> token(match('=') ? TokenType.LessEqual : TokenType.Less);

            case '+' -> token(match('=') ? TokenType.PlusEqual : TokenType.Plus);
            case '-' -> token(match('=') ? TokenType.MinusEqual : TokenType.Minus);
            case '*' -> token(match('=') ? TokenType.StarEqual : TokenType.Star);
            case '/' -> token(match('=') ? TokenType.SlashEqual : TokenType.Slash);
            case '%' -> token(match('=') ? TokenType.PercentageEqual : TokenType.Percentage);

            case '"' -> string();

            default -> error("Unexpected character");
        };
    }

    private boolean isTwice(char current, char c) {
        if (current == c && peek() == c) {
            advance();
            return true;
        }

        return false;
    }

    private Token number() {
        while (isDigit(peek())) advance();

        if (peek() == '.' && isDigit(peekNext())) {
            advance();

            while (isDigit(peek())) advance();

            return token(TokenType.Float);
        }

        return token(TokenType.Int);
    }

    private Token string() {
        while (peek() != '"' && !isAtEnd()) {
            if (peek() == '\n') line++;
            advance();
        }

        if (isAtEnd()) return error("Unterminated string.");

        advance();
        return token(TokenType.String);
    }

    private Token identifier() {
        while (isAlpha(peek()) || isDigit(peek())) advance();
        return token(identifierType());
    }

    private TokenType identifierType() {
        return switch (sb.charAt(0)) {
            case 'n' -> checkKeyword(1, "ull", TokenType.Null);
            case 't' -> checkKeyword(1, "rue", TokenType.True);
            case 'f' -> {
                if (sb.length() > 1) yield switch (sb.charAt(1)) {
                    case 'a' -> checkKeyword(2, "lse", TokenType.False);
                    case 'o' -> checkKeyword(2, "r", TokenType.For);
                    default -> TokenType.Identifier;
                };
                yield TokenType.Identifier;
            }
            case 'i' -> checkKeyword(1, "f", TokenType.If);
            case 'w' -> checkKeyword(1, "hile", TokenType.While);

            default -> TokenType.Identifier;
        };
    }

    private TokenType checkKeyword(int start, String rest, TokenType type) {
        return sb.substring(start).equals(rest) ? type : TokenType.Identifier;
    }

    // Utils

    private boolean match(char expected) {
        if (isAtEnd()) return false;
        if (peek() != expected) return false;

        advance();
        return true;
    }

    private char advance() { return advance(true); }
    private char advance(boolean append) {
        if (append) sb.append(current);

        char prev = current;
        current = next;

        try {
            int i = reader.read();
            next = i == -1 ? '\0' : (char) i;
        } catch (IOException e) {
            e.printStackTrace();
        }

        return prev;
    }

    private char peek() {
        return current;
    }

    private char peekNext() {
        return next;
    }

    private boolean isAtEnd() {
        return peek() == '\0';
    }

    private void skipWhitespace() {
        while (true) {
            char c = peek();

            switch (c) {
                case ' ', '\r', '\t' -> advance(false);
                case '\n' -> {
                    line++;
                    advance(false);
                }
                case '/' -> {
                    if (peekNext() == '/') {
                        while (peek() != '\n' && !isAtEnd()) advance(false);
                    } else return;
                }
                default -> { return; }
            }
        }
    }

    private boolean isDigit(char c) {
        return c >= '0' && c <= '9';
    }

    private boolean isAlpha(char c) {
        return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_';
    }

    private Token token(TokenType type) {
        return new Token(type, sb.toString(), line);
    }

    private Token error(String msg) {
        return new Token(TokenType.Error, msg, line);
    }
}