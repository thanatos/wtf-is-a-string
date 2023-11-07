class Example {
    public static void main(String[] args) {
        System.out.println("This is a string with a surrogate pair:");
        System.out.println("(not a counter-example; this is a well formed Unicode string)");
        String s = "ABC: \ud83d\udca9";
        System.out.println(s.length());
        System.out.println(s);

        String s2 = "ABC: \ud83d";
        System.out.println("This is a string with a lone surrogate:");
        System.out.println("(counter-example; this is NOT a well formed Unicode string)");
        System.out.println(s2.length());
        System.out.println(s2);
    }
}
